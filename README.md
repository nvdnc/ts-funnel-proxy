# ts-funnel-proxy

This is a very minimal HTTP reverse proxy that exposes your upstream over a Tailscale Funnel.

## Usage

To expose an HTTP server running on `localhost:8080` over a Tailscale Funnel, run the following command:

```bash
export TS_AUTHKEY=xxyyzz
ts-funnel-proxy -upstream http://localhost:8080 -ts-hostname example-hostname
```

### Notes
- If not specified, the default Tailscale hostname is `my-funneled-service`.
- Specify `TS_AUTHKEY` as an environment variable to prevent Tailscale from asking for interactive auth. To prompt for interactive auth, the `tsnet` module will log a login URL to the console.


## Alternatives

- [Tailscale Caddy](https://github.com/tailscale/caddy-tailscale)
  - "This plugin is still very experimental." -The README
  - Once this matures we'll likely switch to this, but in the meantime it wasn't quite working for us.


## LICENSE

This project is licensed under the MIT License.

See [LICENSE](./LICENSE) for details.
