import { Button, Group, Modal, PasswordInput, TextInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";

export default function RegisterButtonAndModal() {
  const [opened, { open, close }] = useDisclosure(false);

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
    console.log(values);
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
          />
          <TextInput
            withAsterisk
            label="Name"
            placeholder="What should we call you?"
            key={form.key("name")}
            {...form.getInputProps("name")}
            mt={"xs"}
          />
          <PasswordInput
            withAsterisk
            label="Password"
            key={form.key("password")}
            {...form.getInputProps("password")}
            mt={"xs"}
          />
          <PasswordInput
            withAsterisk
            label="Password Confirm"
            key={form.key("passwordConfirmation")}
            {...form.getInputProps("passwordConfirmation")}
            mt={"xs"}
          />

          <Group justify="flex-end" mt="md">
            <Button type="submit">Register</Button>
          </Group>
        </form>
      </Modal>

      <Button size="md" onClick={open}>
        Register
      </Button>
    </>
  );
}
