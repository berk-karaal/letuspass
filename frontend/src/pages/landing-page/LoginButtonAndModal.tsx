import { authLogin } from "@/api/letuspass";
import { SchemasBadRequestResponse } from "@/api/letuspass.schemas";
import { useAppDispatch } from "@/store/hooks";
import { userLoggedIn } from "@/store/slices/user";
import {
  Button,
  Group,
  Modal,
  PasswordInput,
  Text,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconCheck } from "@tabler/icons-react";
import { useMutation } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function LoginButtonAndModal() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [opened, { open, close }] = useDisclosure(false);

  const [errorText, setErrorText] = useState<string | null>(null);

  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      email: "",
      password: "",
    },

    validate: {
      email: (value) => (value.length > 0 ? null : "Email is required"),
      password: (value) => (value.length > 0 ? null : "Password is required"),
    },
  });

  const handleSubmit = (values: typeof form.values) => {
    setErrorText(null);
    mutation.mutate({ email: values.email, password: values.password });
  };

  const mutation = useMutation({
    mutationFn: authLogin,
    onSuccess: (data) => {
      notifications.show({
        title: "Login Successful",
        message: "You have successfully logged in.",
        color: "green",
        icon: <IconCheck />,
      });
      dispatch(userLoggedIn({ email: data.email, name: data.name }));
      navigate("/app");
    },
    onError: (error) => {
      if (isAxiosError(error)) {
        switch (error.response?.status) {
          case 400:
            let resp = error.response.data as SchemasBadRequestResponse;
            setErrorText(resp.error);
            break;
          default:
            setErrorText("An error occurred. Please try again later.");
        }
      } else {
        setErrorText("An error occurred. Please try again later.");
      }
    },
  });

  return (
    <>
      <Modal opened={opened} onClose={close} title="Login">
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <TextInput
            withAsterisk
            label="Email"
            placeholder="your@email.com"
            key={form.key("email")}
            {...form.getInputProps("email")}
            disabled={mutation.isPending}
          />
          <PasswordInput
            withAsterisk
            label="Password"
            key={form.key("password")}
            {...form.getInputProps("password")}
            disabled={mutation.isPending}
            mt={"xs"}
          />

          <Text c={"red"} mt={"xs"} display={errorText ? "block" : "none"}>
            {errorText}
          </Text>

          <Group justify="flex-end" mt="md">
            <Button type="submit" loading={mutation.isPending}>
              Login
            </Button>
          </Group>
        </form>
      </Modal>

      <Button size="md" onClick={open}>
        Login
      </Button>
    </>
  );
}
