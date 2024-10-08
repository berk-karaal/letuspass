import { listVaultUsers } from "@/api/letuspass";
import { restQueryRetryFunc } from "@/common/queryRetry";
import { useAppSelector } from "@/store/hooks";
import { Box, Group, Loader, Text } from "@mantine/core";
import { IconUserFilled } from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import RemoveUserButtonAndModal from "../RemoveUserButtonAndModal";
import classes from "./styles.module.css";

function VaultUserBox({
  vaultId,
  userId,
  userEmail,
  userPermissions,
}: {
  vaultId: number;
  userId: number;
  userEmail: string;
  userPermissions: string[];
}) {
  const user = useAppSelector((state) => state.user);

  return (
    <Box py={"sm"} px={"xs"} my={"sm"} className={classes.VaultUserBox}>
      <Group gap={"xs"}>
        <IconUserFilled size={"1.75rem"} />
        <Text size="lg">{userEmail}</Text>
        {user.email !== userEmail && (
          <Box ml={"auto"}>
            <RemoveUserButtonAndModal vaultId={vaultId} userId={userId} />
          </Box>
        )}
      </Group>
      <Text mt={"sm"}>
        <b>Permissions:</b> {userPermissions.join(", ")}
      </Text>
    </Box>
  );
}

export default function VaultUsersList({ vaultId }: { vaultId: number }) {
  const vaultUsersQuery = useQuery({
    queryKey: ["vaultUsers", vaultId],
    queryFn: () => listVaultUsers(Number(vaultId)),
    retry: restQueryRetryFunc,
  });

  return (
    <>
      {!vaultUsersQuery.isSuccess ? (
        <Loader color="gray" />
      ) : (
        vaultUsersQuery.data.map((user) => (
          <VaultUserBox
            key={user.email}
            vaultId={Number(vaultId)}
            userId={user.id}
            userEmail={user.email}
            userPermissions={user.permissions}
          />
        ))
      )}{" "}
    </>
  );
}
