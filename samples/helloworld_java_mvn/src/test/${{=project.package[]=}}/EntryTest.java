package ${{=project.package[.]=}}.tests;

import ${{=project.package[.]=}};
import org.junit.Test;
import static org.junit.Assert.assertEquals;

public class MyTests
{

  @Test
  public void shouldRun()
  {
    Entry.main(null);
  }
}
