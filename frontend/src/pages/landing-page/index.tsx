import { LandingNavbar } from "@/components/LandingNavbar";
import { useAppSelector } from "@/store/hooks";
import { Container, Group, Text } from "@mantine/core";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import LoginButtonAndModal from "./LoginButtonAndModal";
import RegisterButtonAndModal from "./RegisterButtonAndModal";
import classes from "./styles.module.css";

function LandingPage() {
  const navigate = useNavigate();
  const user = useAppSelector((state) => state.user);

  useEffect(() => {
    if (user.isAuthenticated) {
      navigate("/app");
    }
  }, [user, navigate]);

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
