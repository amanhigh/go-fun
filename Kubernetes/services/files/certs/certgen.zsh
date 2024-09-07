# Generate Self Signed Certs

# Prompt user for domain name
read -p "Enter the domain name (e.g. example.com, docker): " domain

# Generate a self-signed certificate
openssl req -x509 -newkey rsa:4096 -sha256 -days 365 -nodes \
  -keyout "${domain}.key" -out "${domain}.crt" \
  -subj "/CN=${domain}" \
  -addext "subjectAltName=DNS:${domain},DNS:*.${domain}"

echo "Certificate for ${domain} generated successfully!"