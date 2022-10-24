# RainTree Simulator <!-- omit in toc -->

- [Code Structure](#code-structure)
- [Feature Completeness](#feature-completeness)
- [Test Generation](#test-generation)
  - [Install Dependencies](#install-dependencies)
  - [Generate Test](#generate-test)

The Python scripts in this package are used to simulate RainTree (in Python) in order to understand, visualize and validate the Golang implementation in the main [V1 implementation](https://github.com/pokt-network/pocket).

It uses a [Breadth First Search](https://en.wikipedia.org/wiki/Breadth-first_search) approach to mimic the real network based implementation of RainTree (implemented in Go), and can be considered to be an alternative "validation client" to verify the real P2P implementation.

The primary goal is to use this library to generate unit tests that can be copy-pasted into the main repo.

## Code Structure

```bash
rain-tree-sim/python
├── README.md # This file
├── requirements.txt # Python requirements to run file
├── simulator.py # Utility functions used to simulate RainTree
└── test_generator.py # The entrypoint used by `make p2p_test_generator` to generate RainTree unit tests
```

## Feature Completeness

- [x] Basic RainTree implementation
- [x] Unit Test generation
- [ ] Redundancy Layer
- [ ] Cleanup Layer
- [ ] Dead / partially visible nodes
- [ ] Fuzz testing
- [ ] Multi-simulation evaluation + plotting

## Test Generation

### Install Dependencies

Make sure you install the corresponding dependencies

```bash
    pip3 install -r python/requirements.txt
```

### Generate Test

You can specify 2 parameters to the `p2p_test_generator` make target:

- `rainTreeTestOutputFilename` # the file where the unit test should be written to
- `numRainTreeNodes`: the number of nodes to run in the RainTree simulation

Example:

```bash
rainTreeTestOutputFilename=/tmp/answer.go numRainTreeNodes=12 make p2p_test_generator
```

You can then copy pasta the output from `/tmp/answer.go` to `module_raintree_test.go` to add a new unit test.

_NOTE: You must add comments to the tree visualization component manually._
