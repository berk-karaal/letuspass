import { getCurrentUser } from "@/api/letuspass";
import router from "@/routes";
import { useAppDispatch } from "@/store/hooks";
import { startupComplete, userLoggedIn } from "@/store/slices/user";
import { notifications } from "@mantine/notifications";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useEffect } from "react";
import { RouterProvider } from "react-router-dom";
import { restQueryRetryFunc } from "./common/queryRetry";

function App() {
  const dispatch = useAppDispatch();

  const currentUserQuery = useQuery({
    queryKey: ["current-user"],
    queryFn: getCurrentUser,
    refetchOnWindowFocus: false,
    staleTime: Infinity,
    gcTime: Infinity,
    retry: restQueryRetryFunc,
  });

  useEffect(() => {
    if (currentUserQuery.isError) {
      dispatch(startupComplete());
      if (
        axios.isAxiosError(currentUserQuery.error) &&
        currentUserQuery.error.response?.status !== 401
      ) {
        notifications.show({
          title: "Failed to fetch current user",
          message: "Failed to fetch current user.",
          color: "red",
        });
      }
    } else if (currentUserQuery.data) {
      let savedPrivateKey = localStorage.getItem("privateKey");
      if (!savedPrivateKey) {
        savedPrivateKey = "";
        notifications.show({
          title: "Private key not found",
          message: "Please log-out and log-in again.",
          color: "red",
          autoClose: false,
        });
      }

      dispatch(
        userLoggedIn({
          email: currentUserQuery.data.email,
          name: currentUserQuery.data.name,
          privateKey: savedPrivateKey,
        })
      );
      dispatch(startupComplete());
    }
  }, [currentUserQuery.isFetching, dispatch]);

  return (
    <>
      <RouterProvider router={router} />
    </>
  );
}

export default App;
