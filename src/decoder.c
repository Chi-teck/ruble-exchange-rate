/**
 * @package RER.
 */

#include <stdlib.h>
#include <string.h>
#include "decoder.h"
#include <libxml/parser.h>

// @todo Handle broken XML.
void parse_rates(char *encoded_xml, rate *rate)
{
    rate->total = 0;
    rate->size = 0;
    rate->records = NULL;
    size_t capacity = 0;

    xmlDocPtr document = xmlParseDoc((xmlChar *) encoded_xml);

    xmlNode *root_element = NULL;
    root_element = xmlDocGetRootElement(document);

    if (root_element == NULL) {
        fprintf(stderr, "Could not parse response from CBR.\n");
        exit(EXIT_FAILURE);
    }

    xmlChar *date = xmlGetProp(root_element, (const xmlChar *) "Date");
    strcpy(rate->date, (char *) date);
    xmlFree(date);

    xmlNode *cur_node = NULL;
    unsigned short i = 0;

    for (cur_node = root_element->children; cur_node; cur_node = cur_node->next) {
        if (cur_node->type != XML_ELEMENT_NODE) {
            continue;
        }

        if (i >= capacity) {
            capacity = capacity ? capacity * 2 : 64;
            rate->records = realloc(rate->records, capacity * sizeof(record));
            if (rate->records == NULL) {
                fprintf(stderr, "Out of memory.\n");
                exit(EXIT_FAILURE);
            }
        }

        xmlNode *valute_node = NULL;
        record record;

        for (valute_node = cur_node->children; valute_node; valute_node = valute_node->next) {
            if (valute_node->type != XML_ELEMENT_NODE) {
                continue;
            }

            xmlChar *content = xmlNodeGetContent(valute_node);

            if (xmlStrcmp(valute_node->name, (const xmlChar *) "CharCode") == 0) {
                strcpy(record.code, (char *) content);
            }
            else if (xmlStrcmp(valute_node->name, (const xmlChar *) "Name") == 0) {
                strcpy(record.name, (char *) content);
            }
            else if (xmlStrcmp(valute_node->name, (const xmlChar *) "Nominal") == 0) {
                record.nominal = atoi((char *) content);
            }
            else if (xmlStrcmp(valute_node->name, (const xmlChar *) "Value") == 0) {
                char* comma = strpbrk((char *) content, ",");
                if (comma != NULL) {
                    *comma = '.';
                }
                record.value = atof((char *) content);
            }

            xmlFree(content);
        }

        rate->records[i] = record;
        i++;
    }

    rate->total = i;

    xmlFreeDoc(document);
}
