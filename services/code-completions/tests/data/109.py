#!/usr/bin/env python3
"""
Generate Kubernetes Service configuration for a range of ports
based on cotun-service.yaml template.
"""

import argparse
import yaml
import sys

def parse_port_range(port_range_str):
    """
    Parse port range string into start and end ports.
    
    Args:
        port_range_str (str): Port range string in format "start-end" or single port
        
    Returns:
        tuple: (start_port, end_port)
    """
    if '-' in port_range_str:
        start_port, end_port = port_range_str.split('-')
        start_port = int(start_port.strip())
        end_port = int(end_port.strip())
        
        if start_port <= 0 or end_port <= 0:
            raise ValueError("Port numbers must be positive")
        if start_port > end_port:
            raise ValueError("Start port must be less than or equal to end port")
            
        return start_port, end_port
    else:
        port = int(port_range_str.strip())
        if port <= 0:
            raise ValueError("Port number must be positive")
        return port, port

def main():
    print("Port range parser")
    try:
        start, end = parse_port_range("8080-8090")
        print(f"Port range: {start} to {end}")
    except ValueError as e:
        print(f"Error: {e}")
