import { authLogout } from "@/api/letuspass";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { userLoggedOut } from "@/store/slices/user";
import {
  Box,
  Button,
  Container,
  Group,
  Menu,
  rem,
  Title,
  useComputedColorScheme,
  useMantineColorScheme,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { IconLogout2, IconMoon, IconSun, IconUser } from "@tabler/icons-react";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import classes from "./styles.module.css";

export function AppShellNavbar() {
  const user = useAppSelector((state) => state.user);

  const { setColorScheme } = useMantineColorScheme();
  const computedColorScheme = useComputedColorScheme("light", {
    getInitialValueInEffect: true,
  });

  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const logoutMutation = useMutation({
    mutationFn: authLogout,
    onSuccess: () => {
      notifications.show({
        title: "Logout Successful",
        message: "You have logged-out.",
        color: "green",
      });
      dispatch(userLoggedOut());
      navigate("/");
    },
    onError: () => {
      notifications.show({
        title: "Logout Failed",
        message: "Failed to log-out.",
        color: "red",
      });
    },
  });

  const logoutOnClick = () => {
    logoutMutation.mutate({});
  };

  return (
    <Box className={classes.header} py={"xs"} mb={"sm"}>
      <Container>
        <Group justify="space-between">
          <Title>LetusPass</Title>
          <Menu shadow="md" width={rem(200)} withArrow offset={3}>
            <Menu.Target>
              <Button
                variant="transparent"
                color={computedColorScheme == "light" ? "dark" : "gray"}
                rightSection={<IconUser size={"1.5rem"} />}
                px={"xs"}
              >
                {user.name}
              </Button>
            </Menu.Target>

            <Menu.Dropdown>
              <Menu.Item
                leftSection={
                  computedColorScheme == "light" ? (
                    <IconMoon size={"1.1rem"} />
                  ) : (
                    <IconSun size={"1.1rem"} />
                  )
                }
                onClick={() =>
                  setColorScheme(
                    computedColorScheme === "light" ? "dark" : "light"
                  )
                }
              >
                Change Theme
              </Menu.Item>
              <Menu.Item
                leftSection={<IconLogout2 size={"1.1rem"} />}
                color="red"
                onClick={logoutOnClick}
              >
                Logout
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>
        </Group>
      </Container>
    </Box>
  );
}
