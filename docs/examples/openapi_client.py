#!/usr/bin/env python3
"""
Example: Using AmneziaWG Easy API with OpenAPI-generated client

This example shows how to interact with the API using a Python client
generated from the OpenAPI specification.

Requirements:
    pip install requests

Generate a proper client (optional):
    npm install -g @openapitools/openapi-generator-cli
    openapi-generator-cli generate \
        -i http://localhost:51821/api/openapi.yaml \
        -g python \
        -o ./amnezia-client
"""

import requests
import json
from typing import Optional, List, Dict


class AmneziaWGClient:
    """Simple client for AmneziaWG Easy API"""

    def __init__(self, base_url: str = "http://localhost:51821"):
        self.base_url = base_url
        self.session = requests.Session()

    def login(self, password: str, remember: bool = False) -> bool:
        """
        Login and create a session
        
        Args:
            password: Admin password
            remember: Remember the session
            
        Returns:
            True if login successful
        """
        response = self.session.post(
            f"{self.base_url}/api/session",
            json={"password": password, "remember": remember}
        )
        response.raise_for_status()
        return response.json().get("success", False)

    def logout(self) -> bool:
        """Logout and destroy session"""
        response = self.session.delete(f"{self.base_url}/api/session")
        response.raise_for_status()
        return response.json().get("success", False)

    def get_clients(self) -> List[Dict]:
        """Get all WireGuard clients"""
        response = self.session.get(f"{self.base_url}/api/wireguard/client")
        response.raise_for_status()
        return response.json()

    def create_client(
        self,
        name: str,
        expired_date: Optional[str] = None,
        jc: Optional[str] = None,
        jmin: Optional[str] = None,
        jmax: Optional[str] = None,
        s1: Optional[str] = None,
        s2: Optional[str] = None,
        h1: Optional[str] = None,
        h2: Optional[str] = None,
        h3: Optional[str] = None,
        h4: Optional[str] = None,
    ) -> bool:
        """
        Create a new WireGuard client
        
        Args:
            name: Client name
            expired_date: Expiration date (YYYY-MM-DD)
            jc, jmin, jmax, s1, s2, h1, h2, h3, h4: AmneziaWG parameters
            
        Returns:
            True if creation successful
        """
        data = {"name": name}
        
        if expired_date:
            data["expiredDate"] = expired_date
            
        # Add AmneziaWG parameters if provided
        amnezia_params = {
            "jc": jc, "jmin": jmin, "jmax": jmax,
            "s1": s1, "s2": s2,
            "h1": h1, "h2": h2, "h3": h3, "h4": h4
        }
        
        for key, value in amnezia_params.items():
            if value is not None:
                data[key] = value
        
        response = self.session.post(
            f"{self.base_url}/api/wireguard/client",
            json=data
        )
        response.raise_for_status()
        return response.json().get("success", False)

    def delete_client(self, client_id: str) -> bool:
        """Delete a client"""
        response = self.session.delete(
            f"{self.base_url}/api/wireguard/client/{client_id}"
        )
        response.raise_for_status()
        return response.json().get("success", False)

    def enable_client(self, client_id: str) -> bool:
        """Enable a client"""
        response = self.session.post(
            f"{self.base_url}/api/wireguard/client/{client_id}/enable"
        )
        response.raise_for_status()
        return response.json().get("success", False)

    def disable_client(self, client_id: str) -> bool:
        """Disable a client"""
        response = self.session.post(
            f"{self.base_url}/api/wireguard/client/{client_id}/disable"
        )
        response.raise_for_status()
        return response.json().get("success", False)

    def download_config(self, client_id: str, output_file: str) -> bool:
        """Download client configuration"""
        response = self.session.get(
            f"{self.base_url}/api/wireguard/client/{client_id}/configuration"
        )
        response.raise_for_status()
        
        with open(output_file, "w") as f:
            f.write(response.text)
        
        return True

    def get_qrcode(self, client_id: str, output_file: str) -> bool:
        """Download QR code as SVG"""
        response = self.session.get(
            f"{self.base_url}/api/wireguard/client/{client_id}/qrcode.svg"
        )
        response.raise_for_status()
        
        with open(output_file, "w") as f:
            f.write(response.text)
        
        return True

    def backup_configuration(self, output_file: str) -> bool:
        """Backup WireGuard configuration"""
        response = self.session.get(f"{self.base_url}/api/wireguard/backup")
        response.raise_for_status()
        
        with open(output_file, "w") as f:
            f.write(response.text)
        
        return True

    def restore_configuration(self, backup_file: str) -> bool:
        """Restore WireGuard configuration from backup"""
        with open(backup_file, "r") as f:
            backup_data = f.read()
        
        response = self.session.put(
            f"{self.base_url}/api/wireguard/restore",
            json={"file": backup_data}
        )
        response.raise_for_status()
        return response.json().get("success", False)


def main():
    """Example usage"""
    # Initialize client
    client = AmneziaWGClient("http://localhost:51821")
    
    # Login
    print("Logging in...")
    client.login("your_password", remember=True)
    
    # Create a client
    print("Creating client...")
    client.create_client(
        name="test-client",
        expired_date="2025-12-31",
        jc="7",
        s1="100",
        s2="100"
    )
    
    # List all clients
    print("\nCurrent clients:")
    clients = client.get_clients()
    for c in clients:
        print(f"  - {c['name']} ({c['address']}) - {'enabled' if c['enabled'] else 'disabled'}")
    
    # Download config for first client
    if clients:
        first_client = clients[0]
        print(f"\nDownloading config for {first_client['name']}...")
        client.download_config(first_client['id'], f"{first_client['name']}.conf")
        print(f"  Saved to {first_client['name']}.conf")
        
        # Download QR code
        print(f"Downloading QR code...")
        client.get_qrcode(first_client['id'], f"{first_client['name']}-qr.svg")
        print(f"  Saved to {first_client['name']}-qr.svg")
    
    # Backup configuration
    print("\nBacking up configuration...")
    client.backup_configuration("wg0-backup.json")
    print("  Saved to wg0-backup.json")
    
    # Logout
    print("\nLogging out...")
    client.logout()
    
    print("\n✅ Done!")


if __name__ == "__main__":
    # Example: Run with your password
    # python openapi_client.py
    main()

