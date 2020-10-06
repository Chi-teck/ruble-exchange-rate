/**
 * @package RER.
 */

#include <stdio.h>
#include <getopt.h>
#include <curl/curl.h>
#include <libxml/parser.h>
#include <string.h>

typedef struct {
	char code[4];
	unsigned short int nominal;
	char name[150];
	float value;
} record;

typedef struct {
	size_t size;
	record *records;
	char date[11];
	unsigned short total;
} rate;

void parse_rates(char *encoded_xml, rate *buffer);
