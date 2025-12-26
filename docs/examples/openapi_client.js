#!/usr/bin/env node
/**
 * Example: Using AmneziaWG Easy API with JavaScript/Node.js
 * 
 * This example shows how to interact with the API using axios
 * 
 * Requirements:
 *   npm install axios
 * 
 * Generate a proper TypeScript client (optional):
 *   npm install -g @openapitools/openapi-generator-cli
 *   openapi-generator-cli generate \
 *     -i http://localhost:51821/api/openapi.yaml \
 *     -g typescript-axios \
 *     -o ./amnezia-client-ts
 */

const axios = require('axios');
const fs = require('fs');

class AmneziaWGClient {
    /**
     * Initialize API client
     * @param {string} baseURL - Base URL of the API
     */
    constructor(baseURL = 'http://localhost:51821') {
        this.client = axios.create({
            baseURL,
            withCredentials: true
        });
    }

    /**
     * Login and create session
     * @param {string} password - Admin password
     * @param {boolean} remember - Remember the session
     * @returns {Promise<boolean>}
     */
    async login(password, remember = false) {
        const response = await this.client.post('/api/session', {
            password,
            remember
        });
        return response.data.success;
    }

    /**
     * Logout and destroy session
     * @returns {Promise<boolean>}
     */
    async logout() {
        const response = await this.client.delete('/api/session');
        return response.data.success;
    }

    /**
     * Get all WireGuard clients
     * @returns {Promise<Array>}
     */
    async getClients() {
        const response = await this.client.get('/api/wireguard/client');
        return response.data;
    }

    /**
     * Create a new WireGuard client
     * @param {Object} options - Client options
     * @param {string} options.name - Client name
     * @param {string} [options.expiredDate] - Expiration date (YYYY-MM-DD)
     * @param {string} [options.jc] - Junk packet count
     * @param {string} [options.jmin] - Junk min size
     * @param {string} [options.jmax] - Junk max size
     * @param {string} [options.s1] - Init packet junk size
     * @param {string} [options.s2] - Response packet junk size
     * @param {string} [options.h1] - Init packet magic header
     * @param {string} [options.h2] - Response packet magic header
     * @param {string} [options.h3] - Underload packet magic header
     * @param {string} [options.h4] - Transport packet magic header
     * @returns {Promise<boolean>}
     */
    async createClient(options) {
        const response = await this.client.post('/api/wireguard/client', options);
        return response.data.success;
    }

    /**
     * Delete a client
     * @param {string} clientId - Client UUID
     * @returns {Promise<boolean>}
     */
    async deleteClient(clientId) {
        const response = await this.client.delete(`/api/wireguard/client/${clientId}`);
        return response.data.success;
    }

    /**
     * Enable a client
     * @param {string} clientId - Client UUID
     * @returns {Promise<boolean>}
     */
    async enableClient(clientId) {
        const response = await this.client.post(`/api/wireguard/client/${clientId}/enable`);
        return response.data.success;
    }

    /**
     * Disable a client
     * @param {string} clientId - Client UUID
     * @returns {Promise<boolean>}
     */
    async disableClient(clientId) {
        const response = await this.client.post(`/api/wireguard/client/${clientId}/disable`);
        return response.data.success;
    }

    /**
     * Update client name
     * @param {string} clientId - Client UUID
     * @param {string} name - New name
     * @returns {Promise<boolean>}
     */
    async updateClientName(clientId, name) {
        const response = await this.client.put(
            `/api/wireguard/client/${clientId}/name`,
            { name }
        );
        return response.data.success;
    }

    /**
     * Download client configuration
     * @param {string} clientId - Client UUID
     * @param {string} outputFile - Output file path
     * @returns {Promise<boolean>}
     */
    async downloadConfig(clientId, outputFile) {
        const response = await this.client.get(
            `/api/wireguard/client/${clientId}/configuration`
        );
        fs.writeFileSync(outputFile, response.data);
        return true;
    }

    /**
     * Download QR code as SVG
     * @param {string} clientId - Client UUID
     * @param {string} outputFile - Output file path
     * @returns {Promise<boolean>}
     */
    async getQRCode(clientId, outputFile) {
        const response = await this.client.get(
            `/api/wireguard/client/${clientId}/qrcode.svg`
        );
        fs.writeFileSync(outputFile, response.data);
        return true;
    }

    /**
     * Generate one-time download link
     * @param {string} clientId - Client UUID
     * @returns {Promise<boolean>}
     */
    async generateOneTimeLink(clientId) {
        const response = await this.client.post(
            `/api/wireguard/client/${clientId}/generateOneTimeLink`
        );
        return response.data.success;
    }

    /**
     * Backup WireGuard configuration
     * @param {string} outputFile - Output file path
     * @returns {Promise<boolean>}
     */
    async backupConfiguration(outputFile) {
        const response = await this.client.get('/api/wireguard/backup');
        fs.writeFileSync(outputFile, response.data);
        return true;
    }

    /**
     * Restore WireGuard configuration from backup
     * @param {string} backupFile - Backup file path
     * @returns {Promise<boolean>}
     */
    async restoreConfiguration(backupFile) {
        const backupData = fs.readFileSync(backupFile, 'utf8');
        const response = await this.client.put('/api/wireguard/restore', {
            file: backupData
        });
        return response.data.success;
    }

    /**
     * Get Prometheus metrics in JSON format
     * @returns {Promise<Object>}
     */
    async getMetricsJSON() {
        const response = await this.client.get('/metrics/json');
        return response.data;
    }

    /**
     * Get system information
     * @returns {Promise<Object>}
     */
    async getSystemInfo() {
        const [release, lang, session] = await Promise.all([
            this.client.get('/api/release'),
            this.client.get('/api/lang'),
            this.client.get('/api/session')
        ]);

        return {
            version: release.data,
            language: lang.data,
            session: session.data
        };
    }
}

/**
 * Example usage
 */
async function main() {
    try {
        // Initialize client
        const client = new AmneziaWGClient('http://localhost:51821');

        // Login
        console.log('Logging in...');
        await client.login('your_password', true);

        // Get system info
        console.log('\nSystem information:');
        const systemInfo = await client.getSystemInfo();
        console.log(`  Version: ${systemInfo.version}`);
        console.log(`  Language: ${systemInfo.language}`);
        console.log(`  Authenticated: ${systemInfo.session.authenticated}`);

        // Create a client with custom AmneziaWG parameters
        console.log('\nCreating client...');
        await client.createClient({
            name: 'test-client',
            expiredDate: '2025-12-31',
            jc: '7',
            jmin: '50',
            jmax: '1000',
            s1: '100',
            s2: '100'
        });

        // List all clients
        console.log('\nCurrent clients:');
        const clients = await client.getClients();
        clients.forEach(c => {
            const status = c.enabled ? 'enabled' : 'disabled';
            const connected = c.latestHandshakeAt ? '🟢 connected' : '⚪ disconnected';
            console.log(`  - ${c.name} (${c.address}) - ${status} ${connected}`);
        });

        // Download config for first client
        if (clients.length > 0) {
            const firstClient = clients[0];
            
            console.log(`\nDownloading config for ${firstClient.name}...`);
            await client.downloadConfig(
                firstClient.id,
                `${firstClient.name}.conf`
            );
            console.log(`  ✅ Saved to ${firstClient.name}.conf`);

            // Download QR code
            console.log('Downloading QR code...');
            await client.getQRCode(
                firstClient.id,
                `${firstClient.name}-qr.svg`
            );
            console.log(`  ✅ Saved to ${firstClient.name}-qr.svg`);

            // Generate one-time link (if enabled)
            try {
                console.log('Generating one-time link...');
                await client.generateOneTimeLink(firstClient.id);
                
                // Refresh client data to get the link
                const updatedClients = await client.getClients();
                const updatedClient = updatedClients.find(c => c.id === firstClient.id);
                
                if (updatedClient.oneTimeLink) {
                    console.log(`  ✅ Link: http://localhost:51821/cnf/${updatedClient.oneTimeLink}`);
                }
            } catch (err) {
                console.log('  ⚠️  One-time links not enabled');
            }
        }

        // Backup configuration
        console.log('\nBacking up configuration...');
        await client.backupConfiguration('wg0-backup.json');
        console.log('  ✅ Saved to wg0-backup.json');

        // Get metrics (if enabled)
        try {
            console.log('\nFetching metrics...');
            const metrics = await client.getMetricsJSON();
            console.log(`  Configured peers: ${metrics.wireguard_configured_peers}`);
            console.log(`  Enabled peers: ${metrics.wireguard_enabled_peers}`);
            console.log(`  Connected peers: ${metrics.wireguard_connected_peers}`);
        } catch (err) {
            console.log('  ⚠️  Metrics not enabled');
        }

        // Logout
        console.log('\nLogging out...');
        await client.logout();

        console.log('\n✅ Done!');
    } catch (error) {
        console.error('\n❌ Error:', error.response?.data?.error || error.message);
        process.exit(1);
    }
}

// Run the example
if (require.main === module) {
    main();
}

module.exports = AmneziaWGClient;

