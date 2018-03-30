visualizePrediction[data_, method_] := Module[
  {p, predictionplot, dataplot, xs},
  dataplot =
   ListPlot[List @@@ data, PlotStyle -> Red,
    PlotLegends -> {"Data"}];
  xs = data[[All, 1]];
  p = Predict[data, Method -> method];
  predictionplot = Plot[{
     p[x],
     p[x] + StandardDeviation[p[x, "Distribution"]],
     p[x] - StandardDeviation[p[x, "Distribution"]]
     }, {x, Min[xs] - 1, Max[xs] + 1},
    PlotStyle -> {Blue, Gray, Gray}, Filling -> {2 -> {3}},
    Exclusions -> False, PerformanceGoal -> "Speed",
    PlotLegends -> {"Prediction", "Confidence Interval"}];
  Show[predictionplot, dataplot, PlotLabel -> method, ImageSize -> 250]
  ]
