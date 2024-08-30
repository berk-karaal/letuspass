import { Container } from "@mantine/core";
import { Outlet } from "react-router-dom";
import { AppShellNavbar } from "./Navbar";

function AppShell() {
  return (
    <>
      <AppShellNavbar />
      <Container>
        <Outlet />
      </Container>
    </>
  );
}

export default AppShell;
