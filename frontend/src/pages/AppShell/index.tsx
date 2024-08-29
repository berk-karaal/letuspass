import { Anchor, Container, Group, Title } from "@mantine/core";
import { Link, Outlet } from "react-router-dom";

function AppShell() {
  return (
    <>
      <Container>
        <Group justify="space-between">
          <Title>LetusPass</Title>
          <Group>
            <Anchor component={Link} to={"/"}>
              Landing
            </Anchor>
            <Anchor component={Link} to={"/app"}>
              Home
            </Anchor>
            <Anchor component={Link} to={"#"}>
              User
            </Anchor>
            <Anchor onClick={() => null}>Logout</Anchor>
          </Group>
        </Group>
        <hr />
        <Outlet />
      </Container>
    </>
  );
}

export default AppShell;
