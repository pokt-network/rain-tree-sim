# RainTree simulator

This simulator executes RainTree.

## Not Implemented

The cleanup layer for extra redundancy is not implemented in order to represent the overall coverage.

## Configs

Configs are store in `config.json`.

```json
{
  "NumNodesFirstSimulation":               the # of nodes in the first simulation,
  "NumNodesLastSimulation":                the # of nodes in the last simulation,

  "DeadNodePercentage":                    % of nodes that should be configured to be non-responsive (i.e. not propagating messages),
  "FixedDeadNodes":                        if true, use `FixedDeadNodesIndexArray` in order to specify which nodes should be non-responsive; random otherwise,
  "FixedDeadNodesIndexArray":              if `FixedDeadNodes` is true, this array specifies which nodes should be non-responsive; ignored otherwise,

  "InvertCurve":                           if true, invert the viewership curve; used for doomsday scenarios,
  "FixedViewershipPercentage":             if true, uses `FixedViewershipCurveArray` to specify the partial viewership curve; otherwise it is random,
  "FixedViewershipCurveArray":             if `FixedViewershipPercentage` is true, this array specifies each node's viewership (e.g. [90,80,70] implies 90%, 80%, 70& viewership respectively),

  "RandomizePartialAddressBooks":          if true, shuffle the partial address book after computation,
  "PartialViewershipMedian":               the global median viewership percentage of partial address books,
  "PartialViewershipStdDev":               the global std deviation of partial address books relative to `PartialViewershipMedian`,

  "RedundancyLayerRightOn":                turn on the right side redundancy layer (not the cleanup layer),
  "RedundancyLayerLeftOn":                 turn on the left side redundancy layer (not the cleanup layer),

  "OriginatorIndex":                       index of the initial sender; -1 for random,
  "MaxHotlist":                            the # of nodes in the `hostlist` for the first level (i.e. initial propagation),

  "ShowIndividualNodeSimResult":           if true, show individual node sim results in the `ResultFileOutputName`.json/csv,
  "ShowIndividualNodePartialAddressBooks": if true, show individual node partial address books in the `ResultFileOutputName`.json/csv,
  "ResultFileOutputName":                  the prefix of the .json/.csv output files
}
```

// TODO(discuss): `NumberOfNodes` and `EndingNumberOfNodes` should be changed to `NumberOfNodesInGlobalNetwork` and `NumOriginatorNodesInSimulation`

## Results

results.json

```json
{
  "NumberOfNodes":             how many nodes in the network?,
  "Levels":                    what was the number of levels in the network?,
  "AverageRedundancy":         what was the avg number of messages received?,
  "NonDeadCoveragePercentage": what percentage of the network was hit before the cleanup layer,
  "DeadCount":                 how many dead nodes?,
  "ConsecutiveLevelZeroMatrix": {
      how many consecutive zeroes before the cleanup layer: how many occurences
  },
```

```
SimResult.csv
Nodes,Levels,Comms,Redundancy,Coverage,Missed,LongestMiss
Number of nodes, number of levels, number of (total) communications, avg redundancy, NonDeadCoveragePercentage, how many missed?, longest consecutive miss?

```
