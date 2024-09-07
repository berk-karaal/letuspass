import { retrieveVault } from "@/api/letuspass";
import { restQueryRetryFunc } from "@/common/queryRetry";
import { ActionIcon, Box, Group, Loader, Text, Title } from "@mantine/core";
import {
  IconArrowLeft,
  IconBriefcase2,
  IconFileTime,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router-dom";
import VaultLogList from "./VaultLogList";

export default function VaultLogsPage() {
  const { vaultId } = useParams();
  const navigate = useNavigate();

  const vaultQuery = useQuery({
    queryKey: ["vault", Number(vaultId)],
    queryFn: () => retrieveVault(Number(vaultId)),
    retry: restQueryRetryFunc,
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
      <Group mt={"lg"} gap={"xs"} justify="flex-start">
        <IconFileTime size={"1.5rem"} />
        <Text fz={"h3"} fw={"lighter"}>
          Audit Logs
        </Text>
      </Group>
      <Box mt={"md"}>
        <VaultLogList vaultId={Number(vaultId)} />
      </Box>
    </>
  );
}
