## Background
In the world of bank-to-bank payments, the standard for message formats is ISO20022. This is an XML format, and there are many message types defined by XSDs at https://www.iso20022.org/. At the same time, there is increasing usage of APIs for payments. Hence there is a need to represent ISO20022 messages as OpenAPI Specification. To ensure that the mapping is done correctly, a tool to convert XSD to OpenAPI Spec was needed. **xsd2oas** is that tool.
## Usage
**xsd2oas -in XSDfilename -out yamlFilename**
Reads the input XSD, parses it into internal data structures, then writes it out as OpenAPI.
## Features
xsd2json supports the key XSD features, including:
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
xsd2oas has not been extensively tested. XSD is a rich and complex standard, and there are undoubtedly many XSDs that will break the current version.
