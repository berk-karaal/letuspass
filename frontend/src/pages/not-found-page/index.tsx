import { Anchor, Container, Title } from "@mantine/core";
import { Link } from "react-router-dom";

export default function NotFoundPage() {
  return (
    <Container p={"xl"}>
      <Title order={1}>Page not found.</Title>
      <Anchor component={Link} display={"block"} mt={"sm"} to="/">
        Go back to the home page
      </Anchor>
    </Container>
  );
}
