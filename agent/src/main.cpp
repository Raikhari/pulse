#include <iostream>
#include <string>

#include <curl/curl.h>
#include <nlohmann/json.hpp>

#include <unistd.h>
#include <fstream>
#include <thread>
#include <chrono>
#include <ctime>

using json = nlohmann::json;

double get_cpu_usage() {
	std::ifstream file("/proc/stat");
	std::string cpu;

	long user, nice, system, idle;

	file >> cpu >> user >> nice >> system >> idle;

	long idle1 = idle;
	long total1 = user + nice + system + idle;

	std::this_thread::sleep_for(std::chrono::milliseconds(200));

	file.clear();
	file.seekg(0);

	file >> cpu >> user >> nice >> system >> idle;

	long idle2 = idle;
	long total2 = user + nice + system + idle;

	long total_diff = total2 - total1;
	long idle_diff = idle2 - idle1;

	if (total_diff == 0) return 0.0;

	return (1.0 - (double)idle_diff / total_diff) * 100.0;
}

double get_ram_usage() {
	std::ifstream file("/proc/meminfo");

	std::string label;
	long total_mem = 0;
	long free_mem = 0;

	while (file >> label) {
		if (label == "MemTotal:") {
			file >> total_mem;
		}
		if (label == "MemAvailable:") {
			file >> free_mem;
			break;
		}
	}

	if (total_mem == 0) return 0.0;

	return (1.0 - (double)free_mem / total_mem) * 100.0;
}

double get_uptime() {
	std::ifstream file("/proc/uptime");

	double uptime = 0.0;
	if(!(file >> uptime)){
		return 0.0;
	}

	return uptime;
}

double get_load_average_1m(){
	double load[3];
	getloadavg(load, 3);
	return load[0];
}

bool sendMetrics(const std::string& url, const std::string& jsonData){
	CURL* curl = curl_easy_init();

	if (!curl)
	{
		std::cerr << "Failed to init curl\n";
		return false;
	}

	curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
	curl_easy_setopt(curl, CURLOPT_POST, 1L);
	curl_easy_setopt(curl, CURLOPT_POSTFIELDS, jsonData.c_str());

	struct curl_slist* headers = nullptr;
	headers = curl_slist_append(
			headers,
			"Content-Type: application/json"
			);

	curl_easy_setopt(
			curl,
			CURLOPT_HTTPHEADER,
			headers
			);

	CURLcode res = curl_easy_perform(curl);

	curl_slist_free_all(headers);
	curl_easy_cleanup(curl);

	return res == CURLE_OK;
}


int main(int argc, char* argv[]) {

	while(true){	
		//build url
		std::string serverIP = "localhost";
		if (argc > 1){
			serverIP = argv[1];
		}

		std::string url = "http://" + serverIP + ":8080/metrics";

		//Build JSON payload
		json payload;

		//get hostname
		char hostname[256];
		gethostname(hostname, sizeof(hostname));

		//get values
		double cpu = get_cpu_usage();
		double ram = get_ram_usage();
		double uptime = get_uptime(); 
		double load_avg_1m = get_load_average_1m();

		//populate payload array with values
		payload["hostname"] = hostname;
		payload["cpu"] = cpu;
		payload["ram"] = ram;
		payload["uptime"] = uptime;
		payload["load1"] = load_avg_1m;
		payload["timestamp"] = static_cast<long long>(time(nullptr));

		//dump array into data 'json'string
		std::string data = payload.dump();

		std::cout << "Sending to: " << serverIP << std::endl;

		//send data into function
		bool success = sendMetrics(url, data);
		if (success){
			//std::cout << "Metrics sent successfully from " << hostname << std::endl;
			std::time_t now = std::time(nullptr);
			std::string timeStr = std::ctime(&now);
			timeStr.pop_back();
			std::cout
				<< "[OK " << timeStr << "] "
				<< hostname
				<< " → "
				<< serverIP
				<< " | CPU "
				<< cpu
				<< "% RAM "
				<< ram
				<< "%"
				<< std::endl;
		} else {
			std::cout << "Failed to send metrics" << std::endl;
		}

		std::this_thread::sleep_for(std::chrono::seconds(5));
	}
	return 0;
}
