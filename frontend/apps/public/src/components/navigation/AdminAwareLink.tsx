import {
  forwardRef,
  type AnchorHTMLAttributes,
  type MouseEvent,
} from "react";
import {
  Link,
  NavLink,
  type NavLinkProps,
  type LinkProps,
} from "react-router-dom";

import { useAdminMode } from "../../hooks/useAdminMode";

type AdminAwareLinkProps = LinkProps;

export const AdminAwareLink = forwardRef<HTMLAnchorElement, AdminAwareLinkProps>(
  function AdminAwareLinkComponent({ to, onClick, ...rest }, ref) {
    const { appendModeTo, confirmIfUnsaved } = useAdminMode();

    const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
      if (!confirmIfUnsaved()) {
        event.preventDefault();
        return;
      }
      onClick?.(event);
    };

    return (
      <Link
        {...rest}
        ref={ref}
        to={appendModeTo(to)}
        onClick={handleClick}
      />
    );
  },
);

type AdminAwareNavLinkProps = NavLinkProps;

export const AdminAwareNavLink = forwardRef<
  HTMLAnchorElement,
  AdminAwareNavLinkProps
>(function AdminAwareNavLinkComponent({ to, onClick, ...rest }, ref) {
  const { appendModeTo, confirmIfUnsaved } = useAdminMode();

  const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
    if (!confirmIfUnsaved()) {
      event.preventDefault();
      return;
    }
    onClick?.(event);
  };

  return (
    <NavLink
      {...rest}
      ref={ref}
      to={appendModeTo(to)}
      onClick={handleClick}
    />
  );
});

type AdminAwareAnchorProps = AnchorHTMLAttributes<HTMLAnchorElement>;

export const AdminAwareAnchor = forwardRef<HTMLAnchorElement, AdminAwareAnchorProps>(
  function AdminAwareAnchorComponent({ href, onClick, ...rest }, ref) {
    const { confirmIfUnsaved } = useAdminMode();

    const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
      if (!confirmIfUnsaved()) {
        event.preventDefault();
        return;
      }
      onClick?.(event);
    };

    return <a {...rest} ref={ref} href={href} onClick={handleClick} />;
  },
);
