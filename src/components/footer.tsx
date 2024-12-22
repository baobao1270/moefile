function Footer() {
  return (
    <footer className="text-center text-gray-500 text-sm py-4">
      <p>
        &copy; {new Date(import.meta.env.BUILD_TIMESTAMP).getFullYear()} {import.meta.env.COPYRIGHT_HOLDER} | Powered by {import.meta.env.APP_NAME} {import.meta.env.APP_VERSION}
      </p>
    </footer>
  );
}

export default Footer;
