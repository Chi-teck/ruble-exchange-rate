/**
 * @package RER.
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <curl/curl.h>

struct memory {
    char *data;
    size_t size;
};

static size_t write_buffer(void *chunk, size_t size, size_t nmemb, struct memory *buffer)
{
    size_t real_size = size * nmemb;
    size_t new_size = buffer->size + real_size;

    buffer->data = realloc(buffer->data, new_size + 1);
    if (buffer->data == NULL) {
        fprintf(stderr, "Out of memory.\n");
        exit(EXIT_FAILURE);
    }

    memcpy(&buffer->data[buffer->size], chunk, real_size);
    buffer->data[new_size] = '\0';
    buffer->size += real_size;

    return real_size;
}

char *fetch(char *url)
{
    struct memory buffer;
    buffer.data = malloc(10);
    buffer.data[0] = '\0';
    buffer.size = 0;

    CURL *curl;
    CURLcode result;

    curl_global_init(CURL_GLOBAL_DEFAULT);

    curl = curl_easy_init();

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_buffer);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &buffer);

    result = curl_easy_perform(curl);

    curl_easy_cleanup(curl);
    curl_global_cleanup();

    if (result != CURLE_OK) {
        fprintf(stderr, "Error %s\n", curl_easy_strerror(result));
        exit(EXIT_FAILURE);
    }

    return buffer.data;
}


