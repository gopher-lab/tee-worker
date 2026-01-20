#!/bin/bash
set -e

# Azure THIM configuration for DCsv3 VMs
# Set USE_AZURE_THIM=true to bypass PCCS and use Azure's Global Attestation Cache directly
if [ "${USE_AZURE_THIM}" = "true" ] || [ "${USE_AZURE_THIM}" = "1" ]; then
    echo "[entrypoint] Configuring for Azure THIM (DCsv3 mode)..."
    
    # Update sgx_default_qcnl.conf to use Azure Global Attestation Cache
    THIM_URL="${AZURE_THIM_URL:-https://global.acccache.azure.net/sgx/certification/v4/}"
    
    # Update the pccs_url to point to Azure THIM
    sed -i 's#"pccs_url": *"[^"]*"#"pccs_url": "'"${THIM_URL}"'"#' /etc/sgx_default_qcnl.conf
    
    # Azure THIM uses proper certificates, so we can enable secure cert verification
    # But for compatibility, we'll keep it configurable
    if [ "${THIM_USE_SECURE_CERT}" = "true" ]; then
        sed -i 's#"use_secure_cert": *false#"use_secure_cert": true#' /etc/sgx_default_qcnl.conf
    fi
    
    echo "[entrypoint] PCCS URL set to: ${THIM_URL}"
    echo "[entrypoint] Azure THIM configuration complete"
else
    echo "[entrypoint] Using default PCCS configuration"
fi

# Show current configuration for debugging
if [ "${LOG_LEVEL}" = "debug" ]; then
    echo "[entrypoint] Current sgx_default_qcnl.conf:"
    cat /etc/sgx_default_qcnl.conf
fi

# Execute the main command
exec "$@"
