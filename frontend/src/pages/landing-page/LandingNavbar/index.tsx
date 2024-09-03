import { Box, Container, Group, Title } from "@mantine/core";
import ColorSchemeToggle from "./ColorSchemeToggle";
import classes from "./styles.module.css";

export function LandingNavbar() {
  return (
    <Box className={classes.header} py={"xs"} mb={"sm"}>
      <Container>
        <Group justify="space-between">
          <Title>LetusPass</Title>
          <ColorSchemeToggle />
        </Group>
      </Container>
    </Box>
  );
}
