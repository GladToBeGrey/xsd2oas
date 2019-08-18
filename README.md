## Background
In the world of bank-to-bank payments, the standard for message formats is ISO20022. This is an XML format, and there are many message types defined by XSDs at https://www.iso20022.org/. At the same time, there is increasing usage of APIs for payments. Hence there is a need to represent ISO20022 messages as OpenAPI Specification (Swagger). To ensure that the mapping is done correctly, a tool to convert XSD to OpenAPI Spec was needed. **xsd2oas** is that tool.

## Usage
**xsd2oas -in XSDfilename -out yamlFilename [-mask maskfile -path pathfile -ex examplefile -lic -fixup -all]**
- XSDfilename (mandatory string) is the location of the XSD file to process (in)
- yamlFilename (mandatory string) is the location to write the yaml file (out)
- maskfile (string) allows the user to specify fields to include (in)
- pathfile (string) is the location to write the paths file (out)
- examplefile (string) is the location to write the example JSON file (out)
- template (string) is the location of a file containing a template (in)
- servers (string) is a comma-delimited list of server URLs (in)
- endpoint (string) is the path to the endpoint relative to server URL (in)
- lic prints license information
- fixup fixes a Swagger bug that duplicates all uppercase parameters by Camelcasing
- all includes all elements in the path file (if omitted, only mandatory fields are included)

## What it does
xsd2oas reads the input XSD, parses it into internal data structures, then writes it out as OpenAPI (Swagger) yaml. By default it will only include mandatory fields; if all fields are needed, this can be specified by the **all** flag.

Most payment schemes require only a subset of the full ISO20022. The fields to be included can be specified in a maskfile. The maskfile consists of one or more lines in the following format:
**/FIToFICstmrDrctDbt/DrctDbtTxInf/CdtrAcct/Id/Othr/Id                        # Bacs F06 Originating account number**
The path specifies an XSD element to include in an XPath-like syntax. An optional comment can be appended, introduced with #.

Note that xsd2oas will include all fields that are mandatory according to the XSD, as well as fields specified in the maskfile. This means that the generated spec will be a valid subset of the full spec.

It may be useful to have a list of paths to the elements included in the yaml. This will be generated if the **pathfile** option is specified. A pathfile generated using the **all** option is a useful starting point to edit to create a maskfile. 

An example JSON file that conforms to the specification will be generated if the **examplefile** option is specified. The JSON will be populated with quasi-random data, but each field conforms to the validation rules specified for that field (length, pattern, enumeration etc).

A comma-delimited list of one or more server URLs may be specifed as the **servers** command-line parameter. If supplied these will be placed in the yaml file under the **servers:** section, with one url list entry per server.

The endpoint may be specified using the **endpoint** parameter. It will be placed in the yaml file under the **paths:** section.

The title may be specified using the **title** parameter. It will be placed in the yaml file as the value of info/title.

## Template file
In case the default settings for info, servers, paths etc. are not suitable, they can be completely over-ridden by using a template file. If a file path is provided as the **template** parameter, it completely replaces the yaml file until the **components:** section. The flexibility of the template mechanism is increased by means of substitution strings. If the keyword appears in the template file, it is replaced by the specified value.

Keyword|Value
-------|-----
$TITLE|-title value if provided, else root of XSD filename
$URLS|-servers value if provided, split into a list of - url: servername entries; else -url: https://example.com
$PATH|-endpoint value if provided, else root of XSD filename
$ROOT|**Mandatory** in template file; substituted by the name of the root type of the XSD

See **template.txt** for an example corresponding to the default settings.

## Features
xsd2oas supports the key XSD features, including:
- Mapping of XSD inbuilt types to OAS types
- Use of "$ref" to simplify the OAS schema
- Enforcing field presence via "required": [...]
- Enforcing strict compliance via "additionalProperties": false
- Restrictions on strings (length, pattern, enum)
- Restrictions on numbers (min, max)
- Support for XSD choices via "oneOf"

## Attributes
There is no direct support for attributes in OAS, so the following mapping convention is followed:
- Map to an object type
- The object contains key "value": value of XML text
- The object also contains keys "@Attribname", one per attribute.
Example:
`
IntrBkSttlmAmt: {
   'value': 1234
   '@Ccy': 'GBP'
`
## Version support
xsd2oas generates schema files compatible with OAS Version 3.

## Known limitations
xsd2oas has been tested on several ISO20022 message types and versions. However, XSD is a rich and complex standard, and there are undoubtedly many XSDs that will break the current version.
