import { retrieveVault } from "@/api/letuspass";
import { ActionIcon, Box, Group, Loader, Text, Title } from "@mantine/core";
import {
  IconArrowLeft,
  IconBriefcase2,
  IconUserPlus,
  IconUsersGroup,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useNavigate, useParams } from "react-router-dom";
import VaultUsersList from "./VaultUsersList";

function VaultUsersPage() {
  const { vaultId } = useParams();
  const navigate = useNavigate();

  const vaultQuery = useQuery({
    queryKey: ["vault", vaultId],
    queryFn: () => retrieveVault(Number(vaultId)),
    retry: (failureCount: number, error: Error) => {
      if (failureCount > 2) {
        return false;
      }
      if (axios.isAxiosError(error)) {
        if (
          (error.response?.status ?? 500 >= 400) &&
          (error.response?.status ?? 500 < 500)
        ) {
          return false;
        }
      }
      return true;
    },
  });

  return (
    <>
      <Group gap={"sm"} justify="flex-start">
        <ActionIcon
          variant="transparent"
          color="dark"
          onClick={() => navigate(-1)}
        >
          <IconArrowLeft size={"1.5rem"} />
        </ActionIcon>
        <IconBriefcase2 size={"1.75rem"} />
        <Title order={2} fw={"lighter"}>
          {vaultQuery.isSuccess ? (
            vaultQuery.data?.name
          ) : (
            <Loader color="gray" />
          )}
        </Title>
      </Group>

      <Group mt={"lg"} mb={"sm"} gap={"sm"} justify="flex-start">
        <IconUsersGroup size={"1.5rem"} />
        <Text fz={"h3"} fw={"lighter"}>
          Vault Users
        </Text>
        <Box ml={"auto"} mr={"xs"}>
          <ActionIcon variant="transparent" color="dark" onClick={() => null}>
            <IconUserPlus size={"1.5rem"} />
          </ActionIcon>
        </Box>
      </Group>

      <VaultUsersList vaultId={Number(vaultId)} />
    </>
  );
}

export default VaultUsersPage;
