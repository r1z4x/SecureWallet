from fastapi import APIRouter, Depends, HTTPException, status, Request
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy.orm import Session
from typing import List, Dict, Any
from src.config.database import get_db, get_mongodb
from src.services.auth import get_current_user
from src.models.user import User
from src.config.settings import settings
from src.vulnerabilities.injection.sql_injection import SQLInjectionVulnerabilities
from src.vulnerabilities.injection.nosql_injection import NoSQLInjectionVulnerabilities
from src.vulnerabilities.injection.command_injection import CommandInjectionVulnerabilities
from src.vulnerabilities.authentication.weak_auth import WeakAuthenticationVulnerabilities
from src.vulnerabilities.xss.xss_vulnerabilities import XSSVulnerabilities
from src.vulnerabilities.xxe.xxe_vulnerabilities import XXEVulnerabilities
from src.vulnerabilities.deserialization.pickle_injection import PickleInjectionVulnerabilities

router = APIRouter()
security = HTTPBearer()


# SQL Injection Vulnerabilities
@router.get("/sql-injection/basic/user-search")
async def basic_sql_injection_user_search(
    username: str,
    db: Session = Depends(get_db)
):
    """Basic SQL Injection - User search"""
    if settings.vulnerability_level not in ["basic", "medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = SQLInjectionVulnerabilities(db)
    result = vuln.basic_sql_injection_user_search(username)
    return {"results": result}


@router.get("/sql-injection/medium/user-search")
async def medium_sql_injection_user_search(
    username: str,
    db: Session = Depends(get_db)
):
    """Medium SQL Injection - User search with partial sanitization"""
    if settings.vulnerability_level not in ["medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = SQLInjectionVulnerabilities(db)
    result = vuln.medium_sql_injection_user_search(username)
    return {"results": result}


@router.get("/sql-injection/hard/user-search")
async def hard_sql_injection_user_search(
    username: str,
    db: Session = Depends(get_db)
):
    """Hard SQL Injection - Complex injection scenarios"""
    if settings.vulnerability_level not in ["hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = SQLInjectionVulnerabilities(db)
    result = vuln.hard_sql_injection_user_search(username)
    return {"results": result}


@router.get("/sql-injection/expert/user-search")
async def expert_sql_injection_user_search(
    username: str,
    db: Session = Depends(get_db)
):
    """Expert SQL Injection - Union-based injection"""
    if settings.vulnerability_level != "expert":
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = SQLInjectionVulnerabilities(db)
    result = vuln.expert_sql_injection_user_search(username)
    return {"results": result}


# NoSQL Injection Vulnerabilities
@router.get("/nosql-injection/basic/user-search")
async def basic_nosql_injection_user_search(
    username: str,
    mongodb=Depends(get_mongodb)
):
    """Basic NoSQL Injection - User search"""
    if settings.vulnerability_level not in ["basic", "medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = NoSQLInjectionVulnerabilities(mongodb)
    result = vuln.basic_nosql_injection_user_search(username)
    return {"results": result}


@router.get("/nosql-injection/medium/user-search")
async def medium_nosql_injection_user_search(
    username: str,
    mongodb=Depends(get_mongodb)
):
    """Medium NoSQL Injection - JSON injection"""
    if settings.vulnerability_level not in ["medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = NoSQLInjectionVulnerabilities(mongodb)
    result = vuln.medium_nosql_injection_user_search(username)
    return {"results": result}


@router.get("/nosql-injection/hard/user-search")
async def hard_nosql_injection_user_search(
    username: str,
    mongodb=Depends(get_mongodb)
):
    """Hard NoSQL Injection - Operator injection"""
    if settings.vulnerability_level not in ["hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = NoSQLInjectionVulnerabilities(mongodb)
    result = vuln.hard_nosql_injection_user_search(username)
    return {"results": result}


@router.get("/nosql-injection/expert/user-search")
async def expert_nosql_injection_user_search(
    username: str,
    mongodb=Depends(get_mongodb)
):
    """Expert NoSQL Injection - Complex operator injection"""
    if settings.vulnerability_level != "expert":
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = NoSQLInjectionVulnerabilities(mongodb)
    result = vuln.expert_nosql_injection_user_search(username)
    return {"results": result}


# Command Injection Vulnerabilities
@router.get("/command-injection/basic/ping")
async def basic_command_injection_ping(host: str):
    """Basic Command Injection - Ping command"""
    if settings.vulnerability_level not in ["basic", "medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = CommandInjectionVulnerabilities()
    result = vuln.basic_command_injection_ping(host)
    return {"result": result}


@router.get("/command-injection/medium/ping")
async def medium_command_injection_ping(host: str):
    """Medium Command Injection - Ping with partial sanitization"""
    if settings.vulnerability_level not in ["medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = CommandInjectionVulnerabilities()
    result = vuln.medium_command_injection_ping(host)
    return {"result": result}


@router.get("/command-injection/hard/ping")
async def hard_command_injection_ping(host: str):
    """Hard Command Injection - Complex injection scenarios"""
    if settings.vulnerability_level not in ["hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = CommandInjectionVulnerabilities()
    result = vuln.hard_command_injection_ping(host)
    return {"result": result}


@router.get("/command-injection/expert/ping")
async def expert_command_injection_ping(host: str):
    """Expert Command Injection - Advanced injection techniques"""
    if settings.vulnerability_level != "expert":
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = CommandInjectionVulnerabilities()
    result = vuln.expert_command_injection_ping(host)
    return {"result": result}


# XSS Vulnerabilities
@router.get("/xss/basic/reflected")
async def basic_reflected_xss(user_input: str):
    """Basic Reflected XSS"""
    if settings.vulnerability_level not in ["basic", "medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = XSSVulnerabilities()
    result = vuln.basic_reflected_xss(user_input)
    return {"html": result}


@router.get("/xss/medium/reflected")
async def medium_reflected_xss(user_input: str):
    """Medium Reflected XSS with partial sanitization"""
    if settings.vulnerability_level not in ["medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = XSSVulnerabilities()
    result = vuln.medium_reflected_xss(user_input)
    return {"html": result}


@router.get("/xss/hard/reflected")
async def hard_reflected_xss(user_input: str):
    """Hard Reflected XSS with complex bypass scenarios"""
    if settings.vulnerability_level not in ["hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = XSSVulnerabilities()
    result = vuln.hard_reflected_xss(user_input)
    return {"html": result}


@router.get("/xss/expert/reflected")
async def expert_reflected_xss(user_input: str):
    """Expert Reflected XSS with advanced bypass techniques"""
    if settings.vulnerability_level != "expert":
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    vuln = XSSVulnerabilities()
    result = vuln.expert_reflected_xss(user_input)
    return {"html": result}


# Weak Authentication Vulnerabilities
@router.post("/weak-auth/basic/password-storage")
async def basic_weak_password_storage(request: Request):
    """Basic Weak Password Storage - Plain text"""
    if settings.vulnerability_level not in ["basic", "medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    password = body.get("password", "")
    
    vuln = WeakAuthenticationVulnerabilities()
    result = vuln.basic_weak_password_storage(password)
    return {"stored_password": result}


@router.post("/weak-auth/medium/password-storage")
async def medium_weak_password_storage(request: Request):
    """Medium Weak Password Storage - MD5 hash"""
    if settings.vulnerability_level not in ["medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    password = body.get("password", "")
    
    vuln = WeakAuthenticationVulnerabilities()
    result = vuln.medium_weak_password_storage(password)
    return {"stored_password": result}


@router.post("/weak-auth/hard/password-storage")
async def hard_weak_password_storage(request: Request):
    """Hard Weak Password Storage - SHA1 hash"""
    if settings.vulnerability_level not in ["hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    password = body.get("password", "")
    
    vuln = WeakAuthenticationVulnerabilities()
    result = vuln.hard_weak_password_storage(password)
    return {"stored_password": result}


@router.post("/weak-auth/expert/password-storage")
async def expert_weak_password_storage(request: Request):
    """Expert Weak Password Storage - Base64 encoding"""
    if settings.vulnerability_level != "expert":
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    password = body.get("password", "")
    
    vuln = WeakAuthenticationVulnerabilities()
    result = vuln.expert_weak_password_storage(password)
    return {"stored_password": result}


# XXE Vulnerabilities
@router.post("/xxe/basic/xml-upload")
async def basic_xxe_xml_upload(request: Request):
    """Basic XXE - XML upload endpoint"""
    if settings.vulnerability_level not in ["basic", "medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    xml_content = body.get("xml_content", "")
    
    vuln = XXEVulnerabilities()
    result = vuln.basic_xxe_file_read(xml_content)
    return result


@router.post("/xxe/medium/xml-upload")
async def medium_xxe_xml_upload(request: Request):
    """Medium XXE - XML upload with remote entity support"""
    if settings.vulnerability_level not in ["medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    xml_content = body.get("xml_content", "")
    
    vuln = XXEVulnerabilities()
    result = vuln.medium_xxe_remote_entity(xml_content)
    return result


@router.post("/xxe/hard/xml-upload")
async def hard_xxe_xml_upload(request: Request):
    """Hard XXE - XML upload with parameter entity injection"""
    if settings.vulnerability_level not in ["hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    xml_content = body.get("xml_content", "")
    
    vuln = XXEVulnerabilities()
    result = vuln.hard_xxe_parameter_entity(xml_content)
    return result


@router.post("/xxe/expert/xml-upload")
async def expert_xxe_xml_upload(request: Request):
    """Expert XXE - XML upload with chained attack"""
    if settings.vulnerability_level != "expert":
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    xml_content = body.get("xml_content", "")
    
    vuln = XXEVulnerabilities()
    result = vuln.expert_xxe_chained_attack(xml_content)
    return result


# Pickle Injection Vulnerabilities
@router.post("/pickle-injection/basic/deserialize")
async def basic_pickle_injection(request: Request):
    """Basic Pickle injection - Direct command execution"""
    if settings.vulnerability_level not in ["basic", "medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    data = body.get("data", "")
    
    vuln = PickleInjectionVulnerabilities()
    result = vuln.basic_pickle_injection(data)
    return result


@router.post("/pickle-injection/medium/deserialize")
async def medium_pickle_injection(request: Request):
    """Medium Pickle injection - File system access"""
    if settings.vulnerability_level not in ["medium", "hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    data = body.get("data", "")
    
    vuln = PickleInjectionVulnerabilities()
    result = vuln.medium_pickle_injection(data)
    return result


@router.post("/pickle-injection/hard/deserialize")
async def hard_pickle_injection(request: Request):
    """Hard Pickle injection - Network access"""
    if settings.vulnerability_level not in ["hard", "expert"]:
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    data = body.get("data", "")
    
    vuln = PickleInjectionVulnerabilities()
    result = vuln.hard_pickle_injection(data)
    return result


@router.post("/pickle-injection/expert/deserialize")
async def expert_pickle_injection(request: Request):
    """Expert Pickle injection - Complex payload with encoding"""
    if settings.vulnerability_level != "expert":
        raise HTTPException(status_code=404, detail="Vulnerability not available")
    
    body = await request.json()
    data = body.get("data", "")
    
    vuln = PickleInjectionVulnerabilities()
    result = vuln.expert_pickle_injection(data)
    return result


# Vulnerability Information
@router.get("/info")
async def get_vulnerability_info():
    """Get information about available vulnerabilities"""
    return {
        "vulnerability_level": settings.vulnerability_level,
        "available_vulnerabilities": {
            "basic": [
                "SQL Injection - Direct string concatenation",
                "NoSQL Injection - Direct string injection",
                "Command Injection - Direct command execution",
                "XSS - Direct output without sanitization",
                "Weak Authentication - Plain text password storage",
                "XXE - Direct file read",
                "Pickle Injection - Direct command execution"
            ],
            "medium": [
                "SQL Injection - Partial sanitization bypass",
                "NoSQL Injection - JSON injection",
                "Command Injection - Partial sanitization",
                "XSS - Partial sanitization",
                "Weak Authentication - MD5 hash",
                "XXE - Remote entity inclusion",
                "Pickle Injection - File system access"
            ],
            "hard": [
                "SQL Injection - Complex injection scenarios",
                "NoSQL Injection - Operator injection",
                "Command Injection - Complex injection scenarios",
                "XSS - Complex bypass scenarios",
                "Weak Authentication - SHA1 hash",
                "XXE - Parameter entity injection",
                "Pickle Injection - Network access"
            ],
            "expert": [
                "SQL Injection - Union-based injection",
                "NoSQL Injection - Complex operator injection",
                "Command Injection - Advanced injection techniques",
                "XSS - Advanced bypass techniques",
                "Weak Authentication - Base64 encoding",
                "XXE - Chained attack with multiple techniques",
                "Pickle Injection - Complex payload with encoding"
            ]
        }
    }
