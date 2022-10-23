import { createElement } from "react";

import Table from "./Table";
import LineGraph from "./LineGraph";
import CenteredBox from "../CenteredBox";
import Typography from "@mui/material/Typography";
import ErrorIcon from "@mui/icons-material/ErrorOutline";

export const VisualizationTypes = [Table, LineGraph];
export const VisualizationColors = [
  ["Red", "rgba(255, 99, 132)"],
  ["Green", "rgba(71, 223, 61)"],
  ["Blue", "rgba(68, 155, 245)"],
  ["Purple", "rgba(103, 68, 245)"],
  ["Orange", "rgba(245, 133, 68)"],
];

export const createVisualizationElement = (visualization) => {
  const metadata = JSON.parse(visualization.metadata);

  for (let i = 0; i < VisualizationTypes.length; i++) {
    if (VisualizationTypes[i].name === metadata.name) {
      return createElement(
        VisualizationTypes[i].view,
        {
          visualization,
          visualizationType: VisualizationTypes[i],
          metadata,
        },
        {}
      );
    }
  }

  return (
    <CenteredBox>
      <ErrorIcon sx={{ fontSize: 50, mb: 1.5 }} />
      <Typography variant="h7" sx={{ userSelect: "none", mb: 2 }}>
        Failed to load visualization, invalid type: {metadata.name}
      </Typography>
    </CenteredBox>
  );
};
