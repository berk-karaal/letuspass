import { LandingNavbar } from "@/components/LandingNavbar";
import { Container, Group, Text } from "@mantine/core";
import LoginButtonAndModal from "./LoginButtonAndModal";
import RegisterButtonAndModal from "./RegsiterButtonAndModal";
import classes from "./styles.module.css";

function LandingPage() {
  return (
    <>
      <LandingNavbar />
      <Container>
        <Text className={classes.slogan} my={"xl"}>
          Share your secret credentials with your team in a secure way!
        </Text>

        <Group justify="space-evenly">
          <LoginButtonAndModal />
          <RegisterButtonAndModal />
        </Group>
      </Container>
    </>
  );
}

export default LandingPage;
