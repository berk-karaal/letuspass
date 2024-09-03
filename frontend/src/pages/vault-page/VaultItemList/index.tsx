import { listVaultItems } from "@/api/letuspass";
import {
  ListVaultItemsOrdering,
  ListVaultItemsParams,
} from "@/api/letuspass.schemas";
import { useVaultPermissionsQuery } from "@/hooks/useVaultPermissionsQuery";
import {
  Box,
  CloseButton,
  Group,
  LoadingOverlay,
  Pagination,
  Select,
  Text,
  TextInput,
} from "@mantine/core";
import { useDebouncedValue } from "@mantine/hooks";
import { IconArrowsSort, IconKey, IconSearch } from "@tabler/icons-react";
import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import CreateVaultItemButtonAndModal from "../CreateVaultItemButtonAndModal";
import classes from "./styles.module.css";

function VaultItemBox({
  itemTitle,
  vaultId,
  itemId,
}: {
  itemTitle: string;
  vaultId: number;
  itemId: number;
}) {
  return (
    <Box
      component={Link}
      to={`/app/vault/${vaultId}/item/${itemId}`}
      style={{
        textDecoration: "none",
        color: "inherit",
      }}
    >
      <Box className={classes.VaultItemBox} my={"xs"}>
        <Group px={"lg"} h={"100%"}>
          <IconKey size={"2rem"} />
          <Text size={"1.3rem"}>{itemTitle}</Text>
        </Group>
      </Box>
    </Box>
  );
}

export default function VaultItemList({ vaultId }: { vaultId: number }) {
  const PAGE_SIZE = 5;
  const [activePage, setActivePage] = useState(1);
  const [totalItemCount, setTotalItemCount] = useState(0);
  const [ordering, setOrdering] = useState("title");
  const orderingValues = [
    { value: "title", label: "Title" },
    { value: "created_at", label: "Created at" },
  ];
  const [searchValue, setSearchValue] = useState("");
  const [searchValueDebounced] = useDebouncedValue(searchValue, 200);

  const [queryValues, setQueryValues] = useState({
    pageNumber: activePage,
    pageSize: PAGE_SIZE,
    ordering: ordering as ListVaultItemsOrdering,
    title: searchValueDebounced,
  });
  const vaultsQuery = useQuery({
    queryKey: ["vaultItems", vaultId, queryValues],
    queryFn: () => {
      let params: ListVaultItemsParams = {
        page: queryValues.pageNumber,
        page_size: queryValues.pageSize,
        ordering: queryValues.ordering,
      };
      if (queryValues.title) {
        params = {
          ...params,
          title: queryValues.title,
        };
      }
      return listVaultItems(vaultId, params);
    },
    placeholderData: keepPreviousData,
    gcTime: 0,
  });

  const vaultPermissionsQuery = useVaultPermissionsQuery(Number(vaultId));

  // Everytime vaultsQuery is fetched, update totalItemCount according to the response.
  // This is used to calculate the total number of pages in the pagination.
  useEffect(() => {
    if (vaultsQuery.isSuccess && vaultsQuery.data?.count >= 0) {
      setTotalItemCount(vaultsQuery.data?.count);
    }
  }, [vaultsQuery]);

  useEffect(() => {
    setActivePage(1);
    setQueryValues((prev) => ({
      ...prev,
      pageNumber: 1,
      title: searchValueDebounced,
    }));
  }, [searchValueDebounced]);

  return (
    <>
      <Group
        justify="space-between"
        mt={"md"}
        mb={"xs"}
        gap={"xs"}
        wrap="nowrap"
      >
        <Select
          rightSection={<IconArrowsSort />}
          checkIconPosition="right"
          data={orderingValues}
          value={ordering}
          onChange={(_value, option) => {
            setOrdering(option.value);
            setActivePage(1);
            setQueryValues((prev) => ({
              ...prev,
              pageNumber: 1,
              ordering: option.value as ListVaultItemsOrdering,
            }));
          }}
        />
        {vaultPermissionsQuery.isSuccess &&
          vaultPermissionsQuery.data.includes("manage_items") && (
            <CreateVaultItemButtonAndModal vaultId={vaultId} />
          )}
      </Group>
      <TextInput
        placeholder="Search"
        leftSection={<IconSearch size={"1rem"} />}
        rightSection={
          <CloseButton
            onClick={() => setSearchValue("")}
            style={{ display: searchValue ? undefined : "none" }}
          />
        }
        value={searchValue}
        onChange={(event) => setSearchValue(event.currentTarget.value)}
      />
      <Box pos={"relative"}>
        <LoadingOverlay
          visible={vaultsQuery.isFetching}
          zIndex={1000}
          overlayProps={{ radius: "sm", blur: 2 }}
        />
        {vaultsQuery.data?.results.map((vaultItem) => (
          <VaultItemBox
            key={vaultItem.id}
            itemId={Number(vaultItem.id)}
            itemTitle={vaultItem.title}
            vaultId={vaultId}
          />
        ))}
        {vaultsQuery.isSuccess && vaultsQuery.data?.count === 0 && (
          <Text ta={"center"} my={"lg"}>
            {searchValueDebounced ? "No item found." : "This vault is empty."}
          </Text>
        )}
      </Box>
      <Group mt={"lg"} justify={"center"}>
        <Pagination
          value={activePage}
          onChange={(value) => {
            setActivePage(value);
            setQueryValues((prev) => ({
              ...prev,
              pageNumber: value,
            }));
          }}
          total={Math.ceil(totalItemCount / PAGE_SIZE)}
        />
      </Group>
    </>
  );
}
