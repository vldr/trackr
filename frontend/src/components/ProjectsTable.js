import { useProjects } from "../contexts/ProjectContext";
import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Paper from "@mui/material/Paper";
import CircularProgress from "@mui/material/CircularProgress";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import IconButton from "@mui/material/IconButton";
import MoreVert from "@mui/icons-material/MoreVert";
import CenteredBox from "./CenteredBox";
import LoadingButton from "@mui/lab/LoadingButton";
import NightsStayOutlinedIcon from "@mui/icons-material/NightsStayOutlined";
import Divider from "@mui/material/Divider";

const ProjectsTable = () => {
  const [projects, loading, error] = useProjects();

  if (loading) {
    return (
      <CenteredBox>
        <CircularProgress />
      </CenteredBox>
    );
  }

  if (error) {
  } else {
  }

  return (
    <Container sx={{ mt: 2 }}>
      <Paper sx={{ p: 3, mb: 2 }}>
        <Box
          sx={{
            display: "flex",
            flexDirection: "row",
            mb: 2,
          }}
        >
          <Typography variant="h6" sx={{ flexGrow: 1 }}>
            Projects
          </Typography>
          <LoadingButton
            type="submit"
            variant="contained"
            sx={{ textTransform: "none" }}
            disableElevation
          >
            Create Project
          </LoadingButton>
        </Box>

        <Divider />

        {projects.length ? (
          <Table aria-label="simple table">
            <TableHead>
              <TableRow>
                <TableCell align="left">Name</TableCell>
                <TableCell align="left">Created At</TableCell>
                <TableCell align="right"></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {projects.map((project) => (
                <TableRow
                  key={project.id}
                  sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                >
                  <TableCell align="left">{project.name}</TableCell>
                  <TableCell align="left">{project.createdAt}</TableCell>
                  <TableCell align="right">
                    <IconButton>
                      <MoreVert />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : (
          <Box
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              color: "darkgray",
            }}
          >
            <NightsStayOutlinedIcon sx={{ fontSize: 100, mt: 10, mb: 3 }} />
            <Typography variant="h5" sx={{ mb: 10, userSelect: "none" }}>
              You currently have no projects.
            </Typography>
          </Box>
        )}
      </Paper>
    </Container>
  );
};

export default ProjectsTable;
