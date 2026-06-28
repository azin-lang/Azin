#include <curl/curl.h>
#include <nlohmann/json.hpp>
#include <archive.h>
#include <archive_entry.h>
#include <iomanip>
#include <fstream>
#include <filesystem>
#include <string_view>
class Installer {
  public:
    int apm_download(const std::string_view& repository_url);
};
