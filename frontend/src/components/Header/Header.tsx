// node modules
import React, { useId, useState } from 'react';
import { observer } from 'mobx-react';
// modules
import { ROUTES } from 'routing/registration';
import { colors } from 'utils/color';
import { useStore } from 'modules/Stores';
import { goToForm } from 'modules';
// components
import { KeyboardArrowDown, Close } from '@mui/icons-material';
import MenuIcon from '@mui/icons-material/Menu';
import { AppBar, SwipeableDrawer, IconButton, ListItem, Menu, MenuItem } from '@mui/material';
// styles
import {
  StyledToolbar,
  StyledLogo,
  StyledStack,
  StyledButton,
  StyledButtonText,
  StyledButtonTextDrawer,
  StyledBurger,
  DrawerHeader,
  StyledDrawerList,
  StyledDrawerButton,
  StyledDrawerContactButton,
  StyledButtonTextAnd,
  StyledContactItemPhoneBlock,
  StyledContactBlockTextLink,
  StyledCatalogItemLink
} from './styles';
import { mockedCatalogItems } from 'components/Catalog/mocks';

const headersData: HeadersData[] = [
  {
    label: 'Головна',
    href: ROUTES.HOME,
    variant: 'text',
    color: colors.text.default,
  },
  {
    label: 'Каталог',
    href: ROUTES.CATALOG,
    variant: 'text',
    color: colors.text.default,
  },
  {
    label: 'Про нас',
    href: ROUTES.ABOUT,
    variant: 'text',
    color: colors.text.default,
  },
  {
    label: 'Зв’язатися',
    variant: 'contained',
    color: colors.text.white,
  },
];

const Header = () => {
  const { isOnRoute } = useStore('RoutingStore');
  const [open, setOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<EventTarget & HTMLButtonElement | null>(null);
  const [openDrop, setOpenDrop] = useState(false);

  const drawerWidth = 326;

  const searchParams = new URLSearchParams(window.location.search);
  const titleSearchParam = searchParams.get('title'); 
  
  const getItemColor = (href: string, itemColor: string) => {
    if (href !== ROUTES.HOME && isOnRoute(href)) {
      return colors.text.orange;
    } else if (href === ROUTES.CATALOG && openDrop) {
      return colors.text.orange;
    } else {
      return itemColor;
    }
  }

  const toggleDrawer = () => {
    setOpen(!open);
  };

  const handleOpen = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    setAnchorEl(event.currentTarget);
    setOpenDrop(true);
  };

  const handleClose = (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
    if (e.currentTarget.localName !== "ul") {
      const menu = document.getElementById("simple-menu").children[2] as HTMLElement;
      const menuBoundary = {
        left: menu.offsetLeft,
        top: e.currentTarget.offsetTop + e.currentTarget.offsetHeight,
        right: menu.offsetLeft + menu.offsetHeight,
        bottom: menu.offsetTop + menu.offsetHeight
      };
      if (
        e.clientX >= menuBoundary.left &&
        e.clientX <= menuBoundary.right &&
        e.clientY <= menuBoundary.bottom &&
        e.clientY >= menuBoundary.top
      ) {
        return;
      }
    }

    setOpenDrop(false);
  };

  return (
      <AppBar position="sticky">
        <StyledToolbar>
          <StyledLogo href="/">
            <img src='assets/logo.svg' alt="Emona logo" />
          </StyledLogo>
    
          <StyledStack direction="row" spacing={1}>
            <StyledBurger
              size="large"
              edge="start"
              aria-label="menu"
              sx={{
                mt: 1,
                mr: 2,
              }}
              onClick={toggleDrawer}
            >
              <MenuIcon fontSize="medium" />
            </StyledBurger>

            {headersData.map(({ label, href, variant, color }) => (
              <StyledButton
                {...{
                  key: useId(),
                  color: 'success',
                  variant,
                  size: 'small',
                  labelColor: getItemColor(href || 'null', color),
                  ...(href ? {href} : {onClick: () => goToForm()}),
                  ...(label === 'Каталог' && {
                    'aria-haspopup': 'true',
                    onMouseOver: (e) => handleOpen(e),
                    onMouseLeave: (e) => handleClose(e),
                    endIcon: <KeyboardArrowDown fontSize="large" />
                  })
                }}
              >
                <StyledButtonText>
                  {label}
                </StyledButtonText>
              </StyledButton>
            ))}
          </StyledStack>
        </StyledToolbar>
        <SwipeableDrawer
          sx={{
            width: drawerWidth,
            flexShrink: 0,
            '& .MuiDrawer-paper': {
              width: drawerWidth,
            },
          }}
          onClose={toggleDrawer}
          onOpen={toggleDrawer}
          anchor="right"
          open={open}
        >
          <DrawerHeader>
            <IconButton style={{ color: colors.background.grey }} onClick={toggleDrawer}>
              <Close />
            </IconButton>
          </DrawerHeader>
          <StyledDrawerList>
            {headersData.slice(0, -1).map(({ label, href, color }) => (
              <ListItem key={label} disablePadding>
                <StyledDrawerButton href={ href || '' }>
                  <StyledButtonTextDrawer color={getItemColor(href || 'null', color)}>
                    {label}
                  </StyledButtonTextDrawer>
                </StyledDrawerButton>
              </ListItem>
            ))}
            <StyledDrawerContactButton href={ROUTES.HOME} variant='contained' size="small">
              <StyledButtonTextDrawer color={colors.text.white} textTransform="none">
                Зв’язатися
              </StyledButtonTextDrawer>
            </StyledDrawerContactButton>
            <StyledButtonTextAnd>
              або
            </StyledButtonTextAnd>
            <StyledContactItemPhoneBlock>
              <StyledContactBlockTextLink href="tel:+380444868610" sx={{ color: colors.text.default }}>
                +38 044 486 86 10
              </StyledContactBlockTextLink>
              <StyledContactBlockTextLink href="tel:+380444868596" sx={{ color: colors.text.default }}>
                +38 044 486 85 96
              </StyledContactBlockTextLink>
            </StyledContactItemPhoneBlock>
          </StyledDrawerList>
        </SwipeableDrawer>
        <Menu
          anchorEl={anchorEl}
          open={openDrop}
          MenuListProps={{
            onMouseLeave: handleClose,
          }}
        >
          { mockedCatalogItems.map(item => (
            <MenuItem selected={ item.title.includes(titleSearchParam || 'false') || false }>
              <StyledCatalogItemLink href={`${ROUTES.CATALOG_ITEM}?title=${item.title}`}>{ item.title }</StyledCatalogItemLink>
            </MenuItem>
          ))}
        </Menu>
      </AppBar>
  );
};

export default observer(Header);
