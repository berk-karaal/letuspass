import { Anchor, Container, Title } from "@mantine/core";
import { Link } from "react-router-dom";

function LandingPage() {
  return (
    <>
      <Container>
        <Title>Landing Page</Title>
        <Anchor component={Link} to="/app">
          App
        </Anchor>
        <br />
        <Anchor component={Link} to="/login">
          Login
        </Anchor>
        <br />
        <Anchor component={Link} to="/register">
          Register
        </Anchor>
      </Container>
    </>
  );
}

export default LandingPage;
