# RainTree Simulator <!-- omit in toc -->

- [Code Structure](#code-structure)
- [Feature Completeness](#feature-completeness)
- [Test Generation](#test-generation)
  - [Install Dependencies](#install-dependencies)
  - [Generate Test](#generate-test)

The Python scripts in this package are used to simulate RainTree (in Python) to understand, visualize and validate the Golang implementation in the primary [V1 implementation](https://github.com/pokt-network/pocket).

It uses a [Breadth First Search](https://en.wikipedia.org/wiki/Breadth-first_search) approach to mimic the actual network-based implementation of RainTree (implemented in Go). It can be considered an alternative "validation client" to verify the P2P implementation.

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

Make sure you install the corresponding dependencies.

```bash
    pip3 install -r python/requirements.txt
```

### Using the Python Simulator 

#### Prerequisite
- Python: `v < 3.10.8` 
- pip3

### Setup local Environment
To run the simulator, do the following:
```
git clone https://github.com/pokt-network/rain-tree-sim.git
cd python
pip3 install -r requirements.txt
```

### Run Simulator
Run the following command:
```
rainTreeTestOutputFilename=/tmp/answer.go numRainTreeNodes=12 make p2p_test_generator
```

If you wish to run different scenarios, you can modify two parameters to the `p2p_test_generator` make target:
- `rainTreeTestOutputFilename`: the file where the unit test should be written to
- `numRainTreeNodes`: the number of nodes to run in the RainTree simulation

### See the results
By default, the test stores the results at: `/tmp/answer.go`. 
Through termainal, you can view it with: `cat /tmp/answer.go` 
If you're using VS-Code, you can open it using: `code /tmp/answer.go` 

**Note: if you are going to be running the Pocket Unit test with the results, leave this file open and read the section below.**  

### Running the Pocket Unit Test
You can run the p2p unit test in the [Pocket repo](https://github.com/pokt-network/pocket/blob/main/p2p/module_raintree_test.go), which will break down the messages that is being sent from one node to another from the pre-defined RainTree test as well as the unique test generated from the Python simulator.  

#### Prerequisite
Before you can run the unit test, you will first need to clone the [Pocket repo](https://github.com/pokt-network/pocket) and follow the local env setup steps outlined [here](https://github.com/pokt-network/pocket/blob/main/docs/development/README.md#lfg---development). 

#### Running Unit Test 

1. Open up the output file located: `/tmp/answer.go` 
2. Add comments to the tree visualization(see below)

```
//                                val_1                                                              
//           ┌──────────────────────┴────────────┬────────────────────────────────┐                  
//         val_4                               val_1                            val_7                
//   ┌───────┴────┬─────────┐            ┌───────┴────┬─────────┐         ┌───────┴────┬─────────┐   
// val_6        val_4     val_8        val_3        val_1     val_5     val_9        val_7     val_2 
```

3. **Inside of the pocket repo** - Navigate to the `module_raintree_test.go` file located in the `/p2p/` folder.
4. Copy and paste the modified `answer.go` file below one of the existing `TestRainTreeNetworkComplete"X"` functions in `module_raintree_test.go` to add a new unit test.
5. In the main Pocket dir you can run two commands:
   - `make test_p2p`: runs the entire `module_raintree_test.go` file 
   - `go test -v -count=1 -v ./p2p/... -run <function-name-inside-answers.go file>`: which will run the function specified in the `module_raintree_test.go` file
