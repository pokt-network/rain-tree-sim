This simulator executes raintree WITHOUT the cleanup layer to better represent the coverage before the extra redundancy.

```
config.json
{
"NumberOfNodes": the 'starting' number of nodes in the simulated network - as you can run multiple simulations at once,
"EndingNumberOfNodes": the 'ending' number of nodes in the simulated network - as you can run multiple simulations at once,
"DeadNodePercentage": % of nodes not responding/propagating,
"FixedDeadNodes": if true, use the FixedDeadNodesIndexArray inorder to kill specific node as certain indices,
"FixedDeadNodesIndexArray": if FixedDeadNodes is on, which nodes would you like to kill?,
"RandomizePartialAddressBooks": would you like to randomize the partial address books? If not, though they'll be partial, they'll be in a fixed order,
"ViewershipPercentageFixed": would you like to have a fixed viewership curve? if true, use ViewershipCurveArray below to set specific values,
"ViewershipCurveArray": if ViewershipPercentageFixed=true, set a specific curve of viewership for each node. Ex. [90,80,70] means 90%, 80%, 70& viewership respectively ,
"TargetPartialViewershipPercentage": the global median viewership percentage of partial addressbooks,
"PartialViewershipStdDev": how far do you want the global std deviation to be from the TargetPartialViewershipPercentage median,
"InvertCurve": invert the viewership curve for fun doomsday scenarios,
"RedundancyLayerRightOn": turn on the right side redundancy layer (not the cleanup layer),
"RedundancyLayerLeftOn": turn on the left side redundancy layer (not the cleanup layer),
"MaxHotlist": how many nodes on the 'hotlist' for first level - initial propagation,
"ShowIndividualNodeResults": show individual node results in the results.json/csv,
"ShowIndividualNodePartialAddressBooks": show individual node addrbooks in the results.json/csv,,
"ResultFileOutputName": the prefix of the .json/.csv output files,
"OriginatorIndex": the initial sender index, use -1 for random
}
```

```
results.json
{
  "NumberOfNodes": how many nodes in the network?,
  "Levels": what was the number of levels in the network?,
  "AverageRedundancy": what was the avg number of messages received?,
  "NonDeadCoveragePercentage": what percentage of the network was hit before the cleanup layer,
  "DeadCount": how many dead nodes?,
  "ConsecutiveLevelZeroMatrix": {
   how many consecutive zeroes before the cleanup layer: how many occurences
  },
```

```
results.csv
Nodes,Levels,Comms,Redundancy,Coverage,Missed,LongestMiss
Number of nodes, number of levels, number of (total) communications, avg redundancy, NonDeadCoveragePercentage, how many missed?, longest consecutive miss?

```