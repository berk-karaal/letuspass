import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface User {
  startupComplete: boolean;
  isAuthenticated: boolean;
  email: string;
  name: string;
  privateKey: string;
}

const initialState: User = {
  startupComplete: false,
  isAuthenticated: false,
  email: "",
  name: "",
  privateKey: "",
};

const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    startupComplete(state) {
      state.startupComplete = true;
    },
    userLoggedIn(
      state,
      action: PayloadAction<{ email: string; name: string }>
    ) {
      state.isAuthenticated = true;
      state.email = action.payload.email;
      state.name = action.payload.name;
      state.privateKey = "hebele hübele";
    },
    userLoggedOut(state) {
      state.isAuthenticated = false;
      state.email = "";
      state.name = "";
      state.privateKey = "";
    },
  },
});

export const { startupComplete, userLoggedIn, userLoggedOut } =
  userSlice.actions;

export default userSlice.reducer;
