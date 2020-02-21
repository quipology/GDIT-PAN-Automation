# Purpose
This tool configures a site-to-site VPN between two Palo Alto Firewalls.

## Features
- Ingest YAML file for VPN configuration details
- Configures firewalls via SSH with provided device's management IP
- Concurrent (configures both devices simultaneously)

## Usage
`./prog -f <YAML config filename>` (if no flag passed, 'config.yml' is default)

## Authors
- [Bobby Williams](https://www.linkedin.com/in/bobby-williams-48222450)

## License
[MIT](LICENSE.md)