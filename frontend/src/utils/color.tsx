const palette = {
  white: '#FFFFFF',
  green: '#388E3C',
  greenDark: '#1B5E20',
  grey: '#101010',
  grey2: '#F5F5F5',
  greyLight: '#424242',
  orange: '#FB8C00'
};

export const colors = {
  background: {
    default: palette.grey2,
    white: palette.white,
    green: palette.green,
    greenDark: palette.greenDark,
    grey: palette.grey,
  },
  text: {
    default: palette.greyLight,
    white: palette.white,
    green: palette.green,
    grey: palette.grey2,
    orange: palette.orange,
    greyLight: palette.greyLight,
    greyDark: palette.grey,
  },
  button: {
    default: palette.green,
    carousel: palette.orange,
  }
};

declare global {
  type ThemeColors =
    | 'background.default'
    | 'background.white'
    | 'background.green'
    | 'background.greenDark'
    | 'background.grey'
    | 'text.default'
    | 'text.white'
    | 'text.green'
    | 'text.grey'
    | 'text.orange'
    | 'text.greyLight'
    | 'text.greyDark'
    | 'button.default'
    | 'button.carousel';
}
