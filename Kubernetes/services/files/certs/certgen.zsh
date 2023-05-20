#Generate Self Signed Certs

# Prompt user for domain name
read -p "Enter the domain name (e.g. example.com): " domain

# Generate a self-signed certificate
openssl req -x509 -newkey rsa:4096 -keyout "${domain}.key" -out "${domain}.crt" \
-subj "/CN=${domain}"

echo "Certificate for ${domain} generated successfully!"
