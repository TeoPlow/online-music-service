FROM envoyproxy/envoy:v1.29-latest
COPY envoy.yaml /etc/envoy/envoy.yaml
CMD ["envoy", "-c", "/etc/envoy/envoy.yaml", "--service-node", "api-gateway"]
