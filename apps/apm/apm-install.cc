#include "apm-install.hpp"
#include <iostream>
#include <print>
#include <string>
#include <string_view>
#include <filesystem>
#include <sstream>
#include <iomanip>
#include <fstream>
#include <curl/curl.h>
#include <archive.h>
#include <archive_entry.h>
#include <openssl/evp.h>

using json = nlohmann::json;
namespace fs = std::filesystem;

int Installer::apm_download(const std::string_view& repository_url)
{
     CURL* curl = curl_easy_init();
     if (!curl) {
         return -1;
     };

     fs::path package_dir = fs::path(".apm");
     fs::create_directories(package_dir);

     std::string tar_path = (package_dir / "package.tar.gz").string();
     std::string tar_url = std::string(repository_url) + "/tarball/HEAD";

     FILE* fp = fopen(tar_path.c_str(), "wb");
     if (!fp) {
      curl_easy_cleanup(curl);
      return -1;
     };

     curl_easy_setopt(curl, CURLOPT_URL, tar_url.c_str());
     curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);
     curl_easy_setopt(curl, CURLOPT_USERAGENT, "libcurl-agent/1.0");
     curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, nullptr);
     curl_easy_setopt(curl, CURLOPT_WRITEDATA, fp);
     curl_easy_setopt(curl, CURLOPT_FAILONERROR, 1L);

     CURLcode res = curl_easy_perform(curl);
     if (res != CURLE_OK) {
      long response_code = 0;
      curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &response_code);
      fclose(fp);
      curl_easy_cleanup(curl);
      return -1;
     };

     fclose(fp);
     curl_easy_cleanup(curl);

     struct archive* a = archive_read_new();
     struct archive* ext = archive_write_disk_new();

     archive_read_support_format_tar(a);
     archive_read_support_filter_gzip(a);

     archive_write_disk_set_options(ext,
         ARCHIVE_EXTRACT_TIME |
         ARCHIVE_EXTRACT_PERM |
         ARCHIVE_EXTRACT_ACL);

     int r = archive_read_open_filename(a, tar_path.c_str(), 10240);
     if (r != ARCHIVE_OK) {
         archive_read_free(a);
         archive_write_free(ext);
         return -1;
     };

     struct archive_entry* entry;
     while (archive_read_next_header(a, &entry) == ARCHIVE_OK)
     {
       std::string out = (package_dir / archive_entry_pathname(entry)).string();

       archive_entry_set_pathname(entry, out.c_str());
       r = archive_write_header(ext, entry);
       if (r != ARCHIVE_OK) {
           continue;
       };

       const void* buff;
       size_t size;
       la_int64_t offset;

       while (archive_read_data_block(a, &buff, &size, &offset) == ARCHIVE_OK) {
         archive_write_data_block(ext, buff, size, offset);
       };

       archive_write_finish_entry(ext);
     };

     archive_read_free(a);
     archive_write_free(ext);

     std::error_code ec;
     fs::remove(tar_path, ec);

     return 0;
}
