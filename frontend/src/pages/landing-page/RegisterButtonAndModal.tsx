import { authRegister } from "@/api/letuspass";
import { SchemasBadRequestResponse } from "@/api/letuspass.schemas";
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
import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { useState } from "react";

export default function RegisterButtonAndModal() {
  const [opened, { open, close }] = useDisclosure(false);

  const [errorText, setErrorText] = useState<string | null>(null);

  const registerMutation = useMutation({
    mutationFn: authRegister,
    onSuccess: () => {
      close();
      form.reset();
      notifications.show({
        title: "Registration Successful",
        message: "You can now log-in.",
        color: "green",
      });
    },
    onError: (error) => {
      if (axios.isAxiosError(error)) {
        switch (error.response?.status) {
          case 400:
            const data = error.response.data as SchemasBadRequestResponse;
            setErrorText(data.error);
            break;
          default:
            setErrorText("Failed to register. Please try again later.");
            break;
        }
      } else {
        setErrorText("Failed to register. Please try again later.");
      }
    },
  });

  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      email: "",
      name: "",
      password: "",
      passwordConfirmation: "",
    },

    validate: {
      email: (value) => (value.length > 0 ? null : "Email is required"),
      name: (value) => (value.length > 0 ? null : "Name is required"),
      password: (value) => (value.length > 0 ? null : "Password is required"),
      passwordConfirmation: (value, formValues) =>
        value === formValues.password ? null : "Passwords should match",
    },
  });

  const handleSubmit = (values: typeof form.values) => {
    setErrorText(null);
    registerMutation.mutate({
      email: values.email,
      name: values.name,
      password: values.password,
    });
  };

  return (
    <>
      <Modal opened={opened} onClose={close} title="Register">
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <TextInput
            withAsterisk
            label="Email"
            placeholder="your@email.com"
            key={form.key("email")}
            {...form.getInputProps("email")}
            disabled={registerMutation.isPending}
          />
          <TextInput
            withAsterisk
            label="Name"
            placeholder="What should we call you?"
            key={form.key("name")}
            {...form.getInputProps("name")}
            mt={"xs"}
            disabled={registerMutation.isPending}
          />
          <PasswordInput
            withAsterisk
            label="Password"
            key={form.key("password")}
            {...form.getInputProps("password")}
            mt={"xs"}
            disabled={registerMutation.isPending}
          />
          <PasswordInput
            withAsterisk
            label="Password Confirm"
            key={form.key("passwordConfirmation")}
            {...form.getInputProps("passwordConfirmation")}
            mt={"xs"}
            disabled={registerMutation.isPending}
          />

          <Text c={"red"} mt={"xs"} display={errorText ? "block" : "none"}>
            {errorText}
          </Text>

          <Group justify="flex-end" mt="md">
            <Button type="submit" loading={registerMutation.isPending}>
              Register
            </Button>
          </Group>
        </form>
      </Modal>

      <Button size="md" onClick={open}>
        Register
      </Button>
    </>
  );
}
