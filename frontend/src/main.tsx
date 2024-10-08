import App from "@/App.tsx";
import { store } from "@/store/store";
import { createTheme, MantineProvider } from "@mantine/core";
import "@mantine/core/styles.css";
import { Notifications } from "@mantine/notifications";
import "@mantine/notifications/styles.css";
import { IconEye, IconEyeOff } from "@tabler/icons-react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Provider as ReduxProvider } from "react-redux";

const theme = createTheme({
  components: {
    PasswordInput: {
      defaultProps: {
        visibilityToggleIcon: ({ reveal }: { reveal: boolean }) =>
          reveal ? (
            <IconEyeOff stroke={1.25} size={"1em"} />
          ) : (
            <IconEye stroke={1.25} size={"1em"} />
          ),
      },
    },
  },
});

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <MantineProvider theme={theme}>
      <QueryClientProvider client={queryClient}>
        <ReduxProvider store={store}>
          <Notifications />
          <App />
        </ReduxProvider>
      </QueryClientProvider>
    </MantineProvider>
  </StrictMode>
);
