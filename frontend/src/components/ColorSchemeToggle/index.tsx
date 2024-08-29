import {
  ActionIcon,
  useComputedColorScheme,
  useMantineColorScheme,
} from "@mantine/core";
import { IconMoon, IconSun } from "@tabler/icons-react";
import classes from "./styles.module.css";

function ColorSchemeToggle() {
  const { setColorScheme } = useMantineColorScheme();
  const computedColorScheme = useComputedColorScheme("light", {
    getInitialValueInEffect: true,
  });

  return (
    <ActionIcon
      onClick={() =>
        setColorScheme(computedColorScheme === "light" ? "dark" : "light")
      }
      variant="default"
      size="lg"
    >
      <IconSun
        width={"22px"}
        height={"22px"}
        className={classes.light}
        stroke={1.25}
      />
      <IconMoon
        width={"22px"}
        height={"22px"}
        className={classes.dark}
        stroke={1.25}
      />
    </ActionIcon>
  );
}

export default ColorSchemeToggle;
