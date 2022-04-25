const palette = {
  white: '#FFFFFF',
  green: '#008A2E',
};

export const colors = {
  background: {
    white: palette.white,
    green: palette.green,
  },
};

declare global {
  type ThemeColors =
    | 'background.white'
    | 'background.green';
}
