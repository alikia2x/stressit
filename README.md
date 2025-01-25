# stressit

stressit is a tool designed to apply controlled CPU and memory load on a system. It can be used for testing system performance, stress testing, or simulating high resource usage scenarios.

## Features

- Apply specific CPU load (percentage of total CPU capacity)
- Allocate a specific amount of memory (in GiB)
- Support for absolute CPU usage mode (targeting system-wide CPU usage)
- Cross-platform support (Linux, macOS, Windows)

## Why stressit?

While tools like `stress-ng` are powerful and feature-rich, they can be overly complex for simple stress testing scenarios. stressit was built with simplicity in mind, offering a straightforward way to apply CPU and memory load without the need to navigate through numerous options.

Additionally, on my MacBook Air, stressit has been observed to have a more significant thermal and power impact, compared to `stress-ng`. This makes stressit particularly useful for scenarios where you need to simulate high CPU load and observe the resulting thermal and power behavior.

## Installation

Pre-built binaries are available in the [Releases](https://github.com/alikia2x/stressit/releases) section.

Download the appropriate binary for your platform, give it execute permissions and move it to a location of your choice.
(for example, to a directory already added to the PATH)

### From Source

1. Clone the repository
2. Install Go 1.21+
3. Build the project:
   ```bash
   go build -o dist/stressit .
   ```

## Usage

```bash
stressit [CPU_CORES]
```

Where [CPU_CORES] can be an integer, a floating point number, or not provided.

For example:

```bash
stressit 0.5 # Simulate 50% CPU load
stressit 2   # Simulate 200% CPU load (2 cores)
stressit     # Simulate 800% CPU load on a 8-core CPU
```

or

```bash
./stressit [flags]
```

### Flags

- `-c float`: CPU load to apply (*100% CPU usage). For example, `-c 0.5` will apply 50% CPU load.
- `-m float`: Memory to allocate in GiB (0 to disable). For example, `-m 2.5` will allocate 2.5 GiB of memory.
- `-a`: Enable absolute CPU usage mode (use -c to set target CPU usage). For example, `stressit -a -c 3` will monitor and attempt to maintain total system CPU usage at 300%.

### Examples

1. Apply 275% CPU load:
   ```bash
   stressit -c 2.75
   ```

2. Allocate 3 GiB of memory:
   ```bash
   stressit -m 3
   ```

3. Target 150% system-wide CPU usage:
   ```bash
   stressit -a -c 1.5
   ```

4. Apply 50% CPU load and allocate 1 GiB of memory:
   ```bash
   stressit -c 0.5 -m 1
   ```

### Warning

The `-a` (absolute) mode is not a precise control. PID-based adjustments may not be accurate enough for all scenarios.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.

[GPL V3](LICENSE)
