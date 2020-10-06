/**
 * @package RER.
 */

#include <stdlib.h>
#include <string.h>
#include "decoder.h"
#include <libxml/parser.h>

#define TOTAL_CURRENCIES 34
#define MAX_CURRENCIES (TOTAL_CURRENCIES + 10)

// @todo Handle broken XML.
void parse_rates(char *encoded_xml, rate *rate)
{
    rate->total = 0;
    rate->size = 0;
    rate->records = malloc(MAX_CURRENCIES * sizeof(record));

    static xmlDocPtr document = NULL;
    document = xmlParseDoc((xmlChar *) encoded_xml);

    xmlNode *root_element = NULL;
    root_element = xmlDocGetRootElement(document);

    if (root_element == NULL) {
        fprintf(stderr, "Could not parse response from CBR.\n");
        exit(EXIT_FAILURE);
    }

    strcpy(rate->date, (char *) xmlGetProp(root_element, (const xmlChar *) "Date"));

    xmlNode *cur_node = NULL;
    unsigned short i = 0;

    for (cur_node = root_element->children; cur_node; cur_node = cur_node->next) {
        if (cur_node->type != XML_ELEMENT_NODE) {
            continue;
        }

        if (i > MAX_CURRENCIES) {
          fprintf(stderr, "Too many currencies in the response.\n");
          exit(EXIT_FAILURE);
        }

        xmlNode *valute_node = NULL;
        record record;

        for (valute_node = cur_node->children; valute_node; valute_node = valute_node->next) {
            if (valute_node->type != XML_ELEMENT_NODE) {
                continue;
            }
            if (xmlStrcmp(valute_node->name, (const xmlChar *) "CharCode") == 0) {
                strcpy(record.code, (char *) xmlNodeGetContent(valute_node));
            }
            if (xmlStrcmp(valute_node->name, (const xmlChar *) "Name") == 0) {
                strcpy(record.name, (char *) xmlNodeGetContent(valute_node));
            }
            if (xmlStrcmp(valute_node->name, (const xmlChar *) "Nominal") == 0) {
                record.nominal = atoi((char *) xmlNodeGetContent(valute_node));
            }
            if (xmlStrcmp(valute_node->name, (const xmlChar *) "Value") == 0) {
                char* value = (char *) xmlNodeGetContent(valute_node);
                char* comma;
                comma = strpbrk(value, ",");
                *comma = '.';
                record.value = atof(value);
            }
        }

        rate->records[i] = record;
        i++;
    }

    rate->total = i;
}
