const palette = {
  white: '#FFFFFF',
  green: '#388E3C',
  grey: '#101010',
};

export const colors = {
  background: {
    white: palette.white,
    green: palette.green,
  },
  text: {
    default: palette.grey,
    white: palette.white,
  }
};

declare global {
  type ThemeColors =
    | 'background.white'
    | 'background.green'
    | 'text.default'
    | 'text.white';
}
