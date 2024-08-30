import { retrieveVault } from "@/api/letuspass";
import { ActionIcon, Box, Group, Loader, Title } from "@mantine/core";
import {
  IconArrowLeft,
  IconBriefcase2,
  IconDotsVertical,
  IconFileTime,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useNavigate, useParams } from "react-router-dom";

function VaultPage() {
  const { vaultId } = useParams();
  const navigate = useNavigate();

  const vaultQuery = useQuery({
    queryKey: ["vault", vaultId],
    queryFn: () => retrieveVault(Number(vaultId)),
    gcTime: 0,
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

        <Box style={{ marginLeft: "auto" }}>
          <ActionIcon variant="transparent" color="dark" onClick={() => null}>
            <IconFileTime size={"1.5rem"} />
          </ActionIcon>
          <ActionIcon variant="transparent" color="dark" onClick={() => null}>
            <IconDotsVertical size={"1.5rem"} />
          </ActionIcon>
        </Box>
      </Group>
      <p>Vault Id: {vaultId}</p>
    </>
  );
}

export default VaultPage;
